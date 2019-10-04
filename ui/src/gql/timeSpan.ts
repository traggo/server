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
        }
    }
`;

export const TimeSpans = gql`
    query TimeSpans($cursor: InputCursor) {
        timeSpans(cursor: $cursor) {
            timeSpans {
                id
                start
                end
                tags {
                    key
                    stringValue
                }
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
        }
    }
`;

export const StopTimer = gql`
    mutation StopTimer($id: Int!, $end: Time!) {
        stopTimeSpan(id: $id, end: $end) {
            id
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
        }
    }
`;

export const UpdateTimeSpan = gql`
    mutation UpdateTimeSpan($id: Int!, $start: Time!, $end: Time, $tags: [InputTimeSpanTag!]) {
        updateTimeSpan(id: $id, start: $start, end: $end, tags: $tags) {
            id
            start
            end
            tags {
                key
                stringValue
            }
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
