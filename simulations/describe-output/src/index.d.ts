declare module "main" {
    export function test(): I32;
}

declare module "xtp:test" {
    interface harness {
        assert(name: PTR, value: I64, reason: PTR): void;
        call(func: PTR, input: PTR): PTR;
        time(func: PTR, input: PTR): I64;
        group(name: PTR): void;
        reset(): void;
        mock_input(): PTR;
    }
}