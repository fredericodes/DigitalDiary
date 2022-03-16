namespace ui.Api {
    public static class Route {
        public static readonly string BaseUrl = "http://localhost:7000";

        public static readonly string AuthRoute = "/auth";
        public static readonly string Login = AuthRoute + "/login";
        public static readonly string RegisterUser = AuthRoute + "/register";

        public static readonly string ApiV1Route = "/api/v1";
        public static readonly string ListDiaryJournals = ApiV1Route + "/journal";
        public static readonly string CreateOrUpdateDiaryJournals = ApiV1Route + "/journal";
        public static readonly string AuthPing = ApiV1Route + "/server-auth-ping";
    }
}