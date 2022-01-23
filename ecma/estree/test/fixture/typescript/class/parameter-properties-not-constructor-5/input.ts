class C {
    not_constructor(
        public readonly pur,
        // Also works on AssignmentPattern
        readonly x = 0,
        public y?: number = 0) {}
}