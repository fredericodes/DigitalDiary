namespace ui.Models.Forms {
    public class FormDiaryJournalEntry {
        public string Date { get; }
        public string Content { get; }

        public FormDiaryJournalEntry(string date, string content) {
            Date = date;
            Content = content;
        }
    }
}