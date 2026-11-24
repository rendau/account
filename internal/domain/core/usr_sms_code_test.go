package core

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rendau/dop/adapters/cache/mem"
	smsMock "github.com/rendau/dop/adapters/sms/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rendau/account/internal/adapters/repo"
	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/account/internal/domain/errs"
)

const testPhone = "77011234567"

// repoStub закрывает только то, что нужно отправке смс: номер есть в базе.
type repoStub struct {
	repo.Repo
}

func (repoStub) UsrGetIdForPhone(context.Context, string) (int64, error) {
	return 1, nil
}

func newSmsCodeTestCore() (*St, *smsMock.St) {
	sms := smsMock.New(nil, true)

	r := &St{
		cache:   mem.New(),
		repo:    repoStub{},
		sms:     sms,
		testing: true,
	}
	r.Usr = NewUsr(r)

	return r, sms
}

func checkCode(r *St, phone string, code int) error {
	return r.Usr.CheckPhoneValidatingCode(context.Background(), &entities.PhoneAndSmsCodeSt{
		Phone:   phone,
		SmsCode: code,
	})
}

// wrongCode возвращает код, который заведомо не совпадает с отправленным.
func wrongCode(sent, i int) int {
	return 1000 + (sent-1000+1+i)%9000
}

func TestCheckPhoneValidatingCode_BlocksAfterWrongAttempts(t *testing.T) {
	r, sms := newSmsCodeTestCore()
	ctx := context.Background()

	require.NoError(t, r.Usr.SendPhoneValidatingCode(ctx, testPhone, true))
	code := sms.PullCode()
	require.NotZero(t, code)

	for i := range unValidateWrongCodeLimit - 1 {
		require.ErrorIs(t, checkCode(r, testPhone, wrongCode(code, i)), errs.WrongSmsCode)
	}

	// последняя разрешённая попытка сразу сообщает, что лимит исчерпан
	require.ErrorIs(t, checkCode(r, testPhone, wrongCode(code, 100)), errs.SmsCheckLimitReached)

	// верный код больше не принимается
	require.ErrorIs(t, checkCode(r, testPhone, code), errs.SmsCheckLimitReached)

	// и новую смс запросить нельзя
	require.ErrorIs(t, r.Usr.SendPhoneValidatingCode(ctx, testPhone, true), errs.SmsCheckLimitReached)
	assert.Zero(t, sms.PullCode(), "смс не должна уходить на закрытый номер")
}

func TestCheckPhoneValidatingCode_CorrectCodeWithinLimit(t *testing.T) {
	r, sms := newSmsCodeTestCore()

	require.NoError(t, r.Usr.SendPhoneValidatingCode(context.Background(), testPhone, true))
	code := sms.PullCode()

	for i := range unValidateWrongCodeLimit - 1 {
		require.ErrorIs(t, checkCode(r, testPhone, wrongCode(code, i)), errs.WrongSmsCode)
	}

	require.NoError(t, checkCode(r, testPhone, code))
}

// Параллельный перебор не должен обходить счётчик: сверку с кодом проходят
// только первые unValidateWrongCodeLimit запросов.
func TestCheckPhoneValidatingCode_ParallelGuessing(t *testing.T) {
	r, sms := newSmsCodeTestCore()

	require.NoError(t, r.Usr.SendPhoneValidatingCode(context.Background(), testPhone, true))
	code := sms.PullCode()

	const guesses = 500

	var wrong, blocked, other atomic.Int64
	var wg sync.WaitGroup

	for i := range guesses {
		wg.Add(1) // wg.Go появился в go 1.25, а модуль на 1.24

		go func() {
			defer wg.Done()

			switch checkCode(r, testPhone, wrongCode(code, i)) {
			case errs.WrongSmsCode:
				wrong.Add(1)
			case errs.SmsCheckLimitReached:
				blocked.Add(1)
			default:
				other.Add(1)
			}
		}()
	}
	wg.Wait()

	assert.EqualValues(t, unValidateWrongCodeLimit-1, wrong.Load())
	assert.EqualValues(t, guesses-unValidateWrongCodeLimit+1, blocked.Load())
	assert.Zero(t, other.Load())
}

// Новая смс не обнуляет счётчик неверных попыток.
func TestSendPhoneValidatingCode_KeepsWrongCount(t *testing.T) {
	r, sms := newSmsCodeTestCore()
	cacheKey := fmt.Sprintf(unValidateCacheKeyTmpl, testPhone)

	require.NoError(t, r.cache.SetJsonObj(cacheKey, usrUNValidatingCacheSt{
		Codes:        []int{1111},
		LastSentTime: time.Now().Add(-2 * unValidateSmsSendInterval),
		WrongCount:   unValidateWrongCodeLimit - 1,
	}, unValidateCacheTimeout))

	require.NoError(t, r.Usr.SendPhoneValidatingCode(context.Background(), testPhone, true))
	code := sms.PullCode()

	var cacheValue usrUNValidatingCacheSt
	ok, err := r.cache.GetJsonObj(cacheKey, &cacheValue)
	require.NoError(t, err)
	require.True(t, ok)

	assert.Equal(t, []int{1111, code}, cacheValue.Codes)
	assert.Equal(t, unValidateWrongCodeLimit-1, cacheValue.WrongCount)
}

// Номера со статичным кодом больше нет: без отправленной смс вход закрыт.
func TestCheckPhoneValidatingCode_NoStaticCodePhone(t *testing.T) {
	r, _ := newSmsCodeTestCore()

	require.ErrorIs(t, checkCode(r, "70000000000", 7456), errs.SmsHasNotSentToPhone)
}

// smsFailStub — провайдер, у которого отправка всегда не удаётся.
type smsFailStub struct{}

func (smsFailStub) Send(string, string) bool      { return false }
func (smsFailStub) SendAsync(string, string) bool { return false }

// Пачка одновременных запросов на отправку не должна обходить интервал:
// уходит ровно одна смс, остальные получают «слишком часто».
func TestSendPhoneValidatingCode_ParallelSending(t *testing.T) {
	r, sms := newSmsCodeTestCore()

	const requests = 90 // заглушка смс хранит не больше 100 сообщений

	var sent, tooFrequent, other atomic.Int64
	var wg sync.WaitGroup

	for range requests {
		wg.Add(1) // wg.Go появился в go 1.25, а модуль на 1.24

		go func() {
			defer wg.Done()

			switch r.Usr.SendPhoneValidatingCode(context.Background(), testPhone, true) {
			case nil:
				sent.Add(1)
			case errs.SmsSendTooFrequent:
				tooFrequent.Add(1)
			default:
				other.Add(1)
			}
		}()
	}
	wg.Wait()

	assert.EqualValues(t, 1, sent.Load())
	assert.EqualValues(t, requests-1, tooFrequent.Load())
	assert.Zero(t, other.Load())
	assert.Len(t, sms.PullAll(), 1)
}

// Смс не ушла — бронь снимается целиком: записи нет, повторить можно сразу.
func TestSendPhoneValidatingCode_FailedSendLeavesNoRecord(t *testing.T) {
	r, sms := newSmsCodeTestCore()
	ctx := context.Background()
	cacheKey := fmt.Sprintf(unValidateCacheKeyTmpl, testPhone)

	r.sms = smsFailStub{}
	require.ErrorIs(t, r.Usr.SendPhoneValidatingCode(ctx, testPhone, true), errs.SmsSendFail)

	_, ok, err := r.cache.Get(cacheKey)
	require.NoError(t, err)
	assert.False(t, ok)

	r.sms = sms
	require.NoError(t, r.Usr.SendPhoneValidatingCode(ctx, testPhone, true))
}

// Смс не ушла при уже существующей записи: прежние коды, время отправки
// и счётчик неверных вводов остаются как были.
func TestSendPhoneValidatingCode_FailedSendRestoresRecord(t *testing.T) {
	r, _ := newSmsCodeTestCore()
	cacheKey := fmt.Sprintf(unValidateCacheKeyTmpl, testPhone)
	prevSentTime := time.Now().Add(-2 * unValidateSmsSendInterval).UTC().Truncate(time.Second)

	require.NoError(t, r.cache.SetJsonObj(cacheKey, usrUNValidatingCacheSt{
		Codes:        []int{1111},
		LastSentTime: prevSentTime,
		WrongCount:   2,
	}, unValidateCacheTimeout))

	r.sms = smsFailStub{}
	require.ErrorIs(t, r.Usr.SendPhoneValidatingCode(context.Background(), testPhone, true), errs.SmsSendFail)

	var cacheValue usrUNValidatingCacheSt
	ok, err := r.cache.GetJsonObj(cacheKey, &cacheValue)
	require.NoError(t, err)
	require.True(t, ok)

	assert.Equal(t, []int{1111}, cacheValue.Codes)
	assert.True(t, prevSentTime.Equal(cacheValue.LastSentTime))
	assert.Equal(t, 2, cacheValue.WrongCount)
}
