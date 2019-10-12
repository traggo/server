import {gql} from 'apollo-boost';

export const Trackers = gql`
    query Trackers {
        timers {
            id
            start
            end
            tags {
                key
                stringValue
            }
            oldStart
        }
    }
`;

export const TimeSpans = gql`
    query TimeSpans($cursor: InputCursor) {
        timeSpans(cursor: $cursor) @connection(key: "AllTimeSpans") {
            timeSpans {
                id
                start
                end
                tags {
                    key
                    stringValue
                }
                oldStart
            }
            cursor {
                startId
                offset
                pageSize
            }
        }
    }
`;

export const StartTimer = gql`
    mutation StartTimer($start: Time!, $tags: [InputTimeSpanTag!]) {
        createTimeSpan(start: $start, tags: $tags) {
            id
            start
            end
            tags {
                key
                stringValue
            }
            oldStart
        }
    }
`;

export const StopTimer = gql`
    mutation StopTimer($id: Int!, $end: Time!) {
        stopTimeSpan(id: $id, end: $end) {
            id
            start
            end
            tags {
                key
                stringValue
            }
            oldStart
        }
    }
`;

export const AddTimeSpan = gql`
    mutation AddTimeSpan($start: Time!, $end: Time!, $tags: [InputTimeSpanTag!]) {
        createTimeSpan(start: $start, end: $end, tags: $tags) {
            id
            start
            end
            tags {
                key
                stringValue
            }
            oldStart
        }
    }
`;

export const UpdateTimeSpan = gql`
    mutation UpdateTimeSpan($id: Int!, $start: Time!, $end: Time, $tags: [InputTimeSpanTag!], $oldStart: Time) {
        updateTimeSpan(id: $id, start: $start, end: $end, tags: $tags, oldStart: $oldStart) {
            id
            start
            end
            tags {
                key
                stringValue
            }
            oldStart
        }
    }
`;
export const RemoveTimeSpan = gql`
    mutation RemoveTimeSpan($id: Int!) {
        removeTimeSpan(id: $id) {
            id
        }
    }
`;
