package controllers

const (
	// Page titles
	PT_HOME        = "Home Page"
	PT_LOGIN       = "Login"
	PT_POSTS       = "Posts"
	PT_PREFERENCES = "Preferences"
	PT_REGISTER    = "Register"
	PT_404         = "Page Not Found"

	// Template names
	TEMPLATE_HOME        = "home"
	TEMPLATE_LOGIN       = "login"
	TEMPLATE_POSTS       = "posts"
	TEMPLATE_PREFERENCES = "preferences"
	TEMPLATE_REGISTER    = "register"
	TEMPLATE_404         = "404"

	// URL constants for various routes
	URL_HOME                = "/"
	URL_LOGIN               = "/login"
	URL_POSTS               = "/posts"
	URL_PREFERENCES         = "/preferences"
	URL_REGISTER            = "/register"
	URL_FORGOT_PASSWORD     = "/accounts/forgot-password"
	URL_RESEND_VERIFICATION = "/accounts/resend-verification"

	// Error messages
	ERR_INVALID_FORM_DATA         = "The submitted form data is invalid. Check your input and try again."
	ERR_USER_NOT_FOUND            = `User with that username not found. Maybe you want to <a href="` + URL_REGISTER + `">register</a>?`
	ERR_LOGIN_INVALID_CREDENTIALS = `The credentials you provided are incorrect. Did you <a href="` + URL_FORGOT_PASSWORD + `">forget your password</a>?`
	ERR_ACCOUNT_DISABLED          = `Your account is disabled or banned. You can reach out to support for assistance.`
	ERR_ACCOUNT_UNABLE_TO_LOGIN   = `You cannot log in at this time. Verify your email or contact support. If you misplaced your verification email, you can <a href="` + URL_RESEND_VERIFICATION + `">request a new one</a>.`
	ERR_SESSION_FAILED_TO_CREATE  = "Failed to create session. Please try again later."
	ERR_SESSION_FAILED_TO_SAVE    = "Failed to save session. Please try again later."
)
