export type Diff<T extends string, U extends string> = ({[P in T]: P} & {[P in U]: never} & {[x: string]: never})[T];
export type Omit<T, K extends Extract<keyof T, string>> = Pick<T, Diff<Extract<keyof T, string>, K>>;
