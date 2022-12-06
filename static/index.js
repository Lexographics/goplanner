const login_text = document.querySelector("#login-text")
const email_tab = document.querySelector("#email-tab")

const form_login_register = document.querySelector("#form-login-register")
const button_login_register = document.querySelector("#login-register-btn")


function setLoginState(login) {
    if (login) {
        login_text.textContent = "Login"
        button_login_register.value = "Login"
        email_tab.style.display = "none"
        form_login_register.action = "/login"
    } else {
        login_text.textContent = "Register"
        button_login_register.value = "Register"
        email_tab.style.display = "block"
        form_login_register.action = "/register"
    }
}