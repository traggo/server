export const stripTypename = <T>(value: T): T => {
    if (value === null || value === undefined) {
        return value;
    }

    if (Array.isArray(value)) {
        // tslint:disable-next-line:no-any
        return (value as any).map((x: any) => stripTypename(x));
    }

    if (typeof value !== 'object') {
        return value;
    }

    Object.values(value).forEach(stripTypename);
    if ('__typename' in value) {
        // @ts-ignore
        delete value.__typename;
    }

    return value;
};
