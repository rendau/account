package errs

import (
	"github.com/rendau/dop/dopErrs"
)

const (
	IdRequired                 = dopErrs.Err("id_required")
	IdAlreadyExists            = dopErrs.Err("id_already_exists")
	ApplicationRequired        = dopErrs.Err("application_required")
	CodeRequired               = dopErrs.Err("code_required")
	RoleRequired               = dopErrs.Err("role_required")
	RolesAndPermShouldBeLoaded = dopErrs.Err("roles_and_perm_should_be_loaded")
	BadRole                    = dopErrs.Err("bad_role")
	PhoneRequired              = dopErrs.Err("phone_required")
	BadPhoneFormat             = dopErrs.Err("bad_phone_format")
	PhoneNotExists             = dopErrs.Err("phone_not_exists")
	PhoneExists                = dopErrs.Err("phone_exists")
	SmsSendLimitReached        = dopErrs.Err("sms_send_limit_reached")
	SmsSendTooFrequent         = dopErrs.Err("sms_send_too_frequent")
	SmsSendFail                = dopErrs.Err("sms_send_fail")
	SmsHasNotSentToPhone       = dopErrs.Err("sms_has_not_sent_to_phone")
	WrongSmsCode               = dopErrs.Err("wrong_sms_code")
	NameRequired               = dopErrs.Err("name_required")
	BadName                    = dopErrs.Err("bad_name")
	AppHasNotPermsUrl          = dopErrs.Err("app_has_not_perms_url")
)
