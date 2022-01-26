namespace ui.Models.Forms {
    public class FormLogin {
        public string Email { get; }
        public string Password { get; }

        public FormLogin(string email, string password) {
            Email = email;
            Password = password;
        }
    }
}