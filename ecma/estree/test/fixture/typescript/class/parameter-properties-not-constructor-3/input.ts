class C {
    not_constructor(
        protected po?,
        private pi?: number,
        public readonly pur,
        // Also works on AssignmentPattern
        readonly x = 0,
        public y?: number = 0) {}
}