export const stripTypename = <T>(value: T): T => {
    if (value === null || value === undefined) {
        return value;
    }

    if (Array.isArray(value)) {
        return value.map((x) => stripTypename(x)) as unknown as T;
    }

    if (typeof value !== 'object') {
        return value;
    }

    const result: Record<string, unknown> = {};
    for (const [key, entry] of Object.entries(value)) {
        if (key !== '__typename') {
            result[key] = stripTypename(entry);
        }
    }

    return result as T;
};
