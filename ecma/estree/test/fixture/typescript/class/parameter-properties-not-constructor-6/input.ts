class C {
    not_constructor(
        // Also works on AssignmentPattern
        readonly x = 0,
        public y?: number = 0) {}
}