import moment from 'moment-timezone';

export interface Range {
    from: string;
    to: string;
}

export const findRange = (selection: {range: Range | null; rangeId: number | null}, ranges: Record<number, Range>): Range => {
    if (selection.rangeId !== null) {
        return exclusiveRange(ranges[selection.rangeId]);
    }
    if (selection.range === null) {
        throw new Error('expected rangeId or range to be non null');
    }
    return exclusiveRange(selection.range);
};

export const exclusiveRange = (range: Range) => ({from: range.from, to: range.to});

export function normalizeDate(date: string): string {
    const d = moment(date);
    if (d.isValid()) {
        return d.utc().format();
    } else {
        return date;
    }
}

export function normalizeRangeDateFormat(range: Range): Range {
    return {from: normalizeDate(range.from), to: normalizeDate(range.to)};
}
