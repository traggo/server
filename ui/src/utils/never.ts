export const expectNever = (value: never): never => {
    throw new Error('expected never but was ' + value);
};
