import moment from 'moment';

interface Success {
    success: true;
    value: moment.Moment;
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

export const parseRelativeTime = (value: string, divide: 'endOf' | 'startOf', nowDate = moment()): Success | Failure => {
    if (isValidDate(value)) {
        return success(asDate(value));
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
        return success(time);
    }

    if (value.indexOf('now') !== -1) {
        return failure("'now' must be at the start");
    }

    return failure("Expected valid date or 'now' at index 0");
};

export const success = (value: moment.Moment): Success => {
    return {success: true, value};
};
export const failure = (error: string): Failure => {
    return {success: false, error};
};

const isOperation = (char: string): char is Operation => {
    return Object.values(Operation).indexOf(char) !== -1;
};
const isUnit = (char: string): char is Unit => {
    return Object.values(Unit).indexOf(char) !== -1;
};

export const isValidDate = (value: string, format?: string) => {
    return asDate(value, format).isValid();
};

export const asDate = (value: string, format = 'YYYY-MM-DD HH:mm') => {
    return moment(value, format, true);
};
export const isSameDate = (from: moment.Moment, to?: moment.Moment): boolean => {
    const fromString = from.format('YYYYMMDD');
    return to === undefined || fromString === to.format('YYYYMMDD');
};
