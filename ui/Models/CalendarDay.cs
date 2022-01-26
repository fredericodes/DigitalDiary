using System;
using System.Collections.Generic;

namespace ui.Models {
    public class CalendarDay {
        public int DayNumber { get; set; }
        public DateTime Date { get; set; }
        public bool IsEmpty { get; set; }
    }
}