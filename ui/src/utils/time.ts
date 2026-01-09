import moment from 'moment-timezone';

interface Success {
    success: true;
    preview: moment.Moment;
    localized: string;
    normalized: string;
}

interface Failure {
    success: false;
    error: string;
}

enum Type {
    Operation = 'operation',
    Unit = 'unit',
    Value = 'value',
}

enum Operation {
    Divide = '/',
    Add = '+',
    Substract = '-',
}

enum Unit {
    Year = 'y',
    Month = 'M',
    Week = 'w',
    Day = 'd',
    Hour = 'h',
    Minute = 'm',
    Second = 's',
}

// tslint:disable-next-line:cyclomatic-complexity mccabe-complexity
export const parseRelativeTime = (value: string, divide: 'endOf' | 'startOf', nowDate = moment()): Success | Failure => {
    for (const format of ['YYYY-MM-DD HH:mm', 'YYYY-MM-DD', 'YYYY-MM-DD[T]HH:mm:ssZ']) {
        if (isValidDate(value, format)) {
            const parsed = asDate(value, format);
            if (divide === 'endOf' && format === 'YYYY-MM-DD') {
                parsed.endOf('day');
            }

            if (format === 'YYYY-MM-DD[T]HH:mm:ssZ') {
                const localDate = parsed.clone().local();
                if (
                    (divide === 'startOf' && parsed.isSame(localDate.startOf('day'), 'second')) ||
                    (divide === 'endOf' && parsed.isSame(localDate.endOf('day'), 'second'))
                ) {
                    value = asDate(value).format('YYYY-MM-DD');
                } else {
                    value = asDate(value).format('YYYY-MM-DD HH:mm');
                }
            }
            return success(
                parsed,
                value,
                parsed
                    .clone()
                    .utc()
                    .format()
            );
        }
    }

    if (value.substr(0, 3) === 'now') {
        let time = nowDate;
        let currentIndex = 'now'.length;

        let expectNext = Type.Operation;
        let lastOperation = Operation.Add;
        let lastNumber = -1;
        while (currentIndex < value.length) {
            const currentChar = value.charAt(currentIndex);
            switch (expectNext) {
                case Type.Operation:
                    if (!isOperation(currentChar)) {
                        return failure('Expected one of / + - at index ' + currentIndex + ' but was ' + currentChar);
                    }
                    lastOperation = currentChar;
                    expectNext = currentChar === Operation.Divide ? Type.Unit : Type.Value;
                    currentIndex++;
                    break;
                case Type.Value:
                    if (isNaN(parseInt(currentChar, 10))) {
                        return failure('Expected number at index ' + currentIndex + ' but was ' + currentChar);
                    }
                    let valueIndex = currentIndex;
                    while (!isNaN(parseInt(value.charAt(valueIndex + 1), 10)) && valueIndex + 1 < value.length) {
                        valueIndex++;
                    }
                    lastNumber = parseInt(value.substr(currentIndex, valueIndex), 10);

                    expectNext = Type.Unit;
                    currentIndex = valueIndex + 1;
                    break;
                case Type.Unit:
                    if (!isUnit(currentChar)) {
                        return failure(
                            'Expected unit (' + Object.values(Unit) + ') at index ' + currentIndex + ' but was ' + currentChar
                        );
                    }

                    // tslint:disable-next-line:no-nested-switch
                    switch (lastOperation) {
                        case Operation.Divide:
                            time = time[divide](currentChar);
                            expectNext = Type.Operation;
                            break;
                        case Operation.Add:
                            time = time.add(lastNumber, currentChar);
                            expectNext = Type.Operation;
                            break;
                        case Operation.Substract:
                            time = time.subtract(lastNumber, currentChar);
                            expectNext = Type.Operation;
                            break;
                        default:
                            throw new Error('oops');
                    }

                    currentIndex++;
                    break;
                default:
                    throw new Error('oopsie');
            }
        }
        if (expectNext === Type.Unit) {
            return failure('Expected unit at the end but got nothing');
        }
        if (expectNext === Type.Value) {
            return failure('Expected number at the end but got nothing');
        }
        return success(time, value, value);
    }

    if (value.indexOf('now') !== -1) {
        return failure("'now' must be at the start");
    }

    return failure("Expected valid date (e.g. 2020-01-01 16:30) or 'now' at index 0");
};

export const success = (value: moment.Moment, localized: string, normalized: string): Success => {
    return {success: true, preview: value, normalized, localized};
};
export const failure = (error: string): Failure => {
    return {success: false, error};
};

const isOperation = (char: string): char is Operation => {
    return Object.values(Operation).indexOf(char as Operation) !== -1;
};
const isUnit = (char: string): char is Unit => {
    return Object.values(Unit).indexOf(char as Unit) !== -1;
};

export const isValidDate = (value: string, format?: string) => {
    return asDate(value, format).isValid();
};

export const asDate = (value: string, format = 'YYYY-MM-DD[T]HH:mm:ssZ') => {
    return moment(value, format, true);
};
export const isSameDate = (from: moment.Moment, to?: moment.Moment): boolean => {
    const fromString = from.format('YYYYMMDD');
    return to === undefined || fromString === to.format('YYYYMMDD');
};
