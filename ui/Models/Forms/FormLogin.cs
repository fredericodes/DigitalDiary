namespace ui.Models.Forms {
    public class FormLogin {
        public string Username { get; }
        public string Password { get; }

        public FormLogin(string username, string password) {
            Username = username;
            Password = password;
        }
    }
}