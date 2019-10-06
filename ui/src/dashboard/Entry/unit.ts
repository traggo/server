interface Unit {
    toUnit: (seconds: number) => number;
    short: string;
    name: string;
}

const Minutes: Unit = {
    toUnit: (s) => s / 60,
    name: 'minutes',
    short: 'm',
};
const Hours: Unit = {
    toUnit: (s) => Minutes.toUnit(s) / 60,
    name: 'hours',
    short: 'h',
};
const Days: Unit = {
    toUnit: (s) => Hours.toUnit(s) / 24,
    name: 'days',
    short: 'd',
};

export const ofSeconds = (seconds: number): Unit => {
    if (seconds < 2 * 60 * 60) {
        return Minutes;
    }

    if (seconds < 30 * 60 * 60) {
        return Hours;
    }

    return Days;
};
