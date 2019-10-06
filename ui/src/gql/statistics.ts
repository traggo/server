import {gql} from 'apollo-boost';

export const Stats = gql`
    query Stats($ranges: [Range!], $tags: [String!], $excludeTags: [InputTimeSpanTag!], $requireTags: [InputTimeSpanTag!]) {
        stats(ranges: $ranges, tags: $tags, excludeTags: $excludeTags, requireTags: $requireTags) {
            start
            end
            entries {
                key
                stringValue
                timeSpendInSeconds
            }
        }
    }
`;
export const Stats2 = gql`
    query Stats2($now: Time!, $stats: InputStatsSelection!) {
        stats: stats2(now: $now, stats: $stats) {
            start
            end
            entries {
                key
                stringValue
                timeSpendInSeconds
            }
        }
    }
`;
