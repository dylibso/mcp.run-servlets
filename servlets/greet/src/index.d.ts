declare module "main" {
  export function call(): I32;
  export function describe(): I32;
}

declare module "extism:host" {
  interface user {
    config_get(ptr: I64): I64;
  }
}
