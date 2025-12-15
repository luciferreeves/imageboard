package config

const (
	// Page titles
	PT_HOME            = "Home Page"
	PT_LOGIN           = "Login"
	PT_FORGOT_PASSWORD = "Forgot Password"
	PT_FORGOT_USERNAME = "Forgot Username"
	PT_POST_LIST       = "Posts"
	PT_POST_NEW        = "Upload New Post"
	PT_POST_SINGLE     = "Post"
	PT_POST_EDIT       = "Edit Post"
	PT_PREFERENCES     = "Preferences"
	PT_REGISTER        = "Register"
	PT_404             = "Page Not Found"
	PT_VERIFY_EMAIL    = "Verify Email"

	// Template names
	TEMPLATE_HOME         = "home"
	TEMPLATE_LOGIN        = "account/login"
	TEMPLATE_REGISTER     = "account/register"
	TEMPLATE_FORGOT       = "account/forgot"
	TEMPLATE_VERIFY_EMAIL = "account/verify_email"
	TEMPLATE_POST_LIST    = "posts/list"
	TEMPLATE_POST_NEW     = "posts/new"
	TEMPLATE_POST_SINGLE  = "posts/single"
	TEMPLATE_POST_EDIT    = "posts/edit"
	TEMPLATE_PREFERENCES  = "preferences"
	TEMPLATE_ERROR        = "error"

	// URL constants for various routes
	URL_HOME                = "/"
	URL_LOGIN               = "/account/login"
	URL_LOGOUT              = "/account/logout"
	URL_REGISTER            = "/account/register"
	URL_VERIFY_EMAIL        = "/account/verify"
	URL_FORGOT_PASSWORD     = "/account/forgot-password"
	URL_RESEND_VERIFICATION = "/account/resend-verification"
	URL_POST_LIST           = "/posts"
	URL_POST_NEW            = "/posts/new"
	URL_PREFERENCES         = "/preferences"

	// Error messages
	ERR_INVALID_FORM_DATA                     = "The submitted form data is invalid. Check your input and try again."
	ERR_USER_NOT_FOUND                        = `User with that username not found. Maybe you want to <a href="` + URL_REGISTER + `">register</a>?`
	ERR_LOGIN_INVALID_CREDENTIALS             = `The credentials you provided are incorrect. Did you <a href="` + URL_FORGOT_PASSWORD + `">forget your password</a>?`
	ERR_ACCOUNT_DISABLED                      = `Your account is disabled or banned. You can reach out to support for assistance.`
	ERR_ACCOUNT_UNABLE_TO_LOGIN               = `You cannot log in at this time. Verify your email or contact support. If you misplaced your verification email, you can <a href="` + URL_RESEND_VERIFICATION + `">request a new one</a>.`
	ERR_PASSWORD_MISMATCH                     = "Entered passwords do not match. Ensure both fields are identical."
	ERR_SESSION_FAILED_TO_CREATE              = "Server failed to create a session. If this issue persists, contact support."
	ERR_SESSION_FAILED_TO_SAVE                = "Server failed to save session data. If this issue persists, contact support."
	ERR_REGISTER_FAILED_TO_CREATE_USER        = "Failed to create user account: "
	ERR_REGISTER_USER_CREATED_EMAIL_FAILED    = "User account created, but failed to send verification email."
	ERR_VERIFY_EMAIL_MISSING_TOKEN            = `Verification token is missing. Check the link you clicked or request a <a href="` + URL_RESEND_VERIFICATION + `">new verification email</a>.`
	ERR_VERIFY_EMAIL_INVALID_OR_EXPIRED_TOKEN = `The verification token is either invalid or has expired. Try requesting a <a href="` + URL_RESEND_VERIFICATION + `">new verification email</a>.`
	ERR_VERIFY_EMAIL_USER_NOT_FOUND           = `User not found for the provided verification token. If you think this is an error, contact support.`
	ERR_VERIFY_EMAIL_ACTIVATION_FAILED        = `Failed to activate your account. If this issue persists, contact support.`

	// Success messages
	SUCCESS_USER_REGISTERED            = "Your account has been created successfully. A verification email has been sent to your email address. You will only be able to log in after verifying your email. If you did not receive the email, you can <a href=\"" + URL_RESEND_VERIFICATION + "\">request a new one</a>."
	SUCCESS_VERIFY_EMAIL               = `Your email has been successfully verified. You can now <a href="` + URL_LOGIN + `">log in</a> to your account.`
	SUCCESS_FORGOT_USERNAME_EMAIL_SENT = "An email has been sent to your email address with all your associated usernames."

	// Non Existent User
	ERR_NO_ACCOUNT_ASSOCIATED_WITH_EMAIL = "No account is associated with the provided email address. Check for typos or consider registering for a new account."
)
