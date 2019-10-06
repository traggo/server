export const stripTypename = <T>(value: T): T => {
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
