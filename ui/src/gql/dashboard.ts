import {gql} from 'apollo-boost';

export const Dashboards = gql`
    query Dashboards {
        dashboards {
            id
            name
            ranges {
                id
                name
                editable
                range {
                    from
                    to
                }
            }
            items {
                id
                title
                total
                entryType
                statsSelection {
                    range {
                        from
                        to
                    }
                    rangeId
                    interval
                    tags
                    excludeTags {
                        key
                        value
                    }
                    includeTags {
                        key
                        value
                    }
                }
                pos {
                    mobile {
                        w
                        h
                        x
                        y
                        minH
                        minW
                    }
                    desktop {
                        w
                        h
                        x
                        y
                        minH
                        minW
                    }
                }
            }
        }
    }
`;

export const RemoveDashboard = gql`
    mutation RemoveDashboard($id: Int!) {
        removeDashboard(id: $id) {
            id
        }
    }
`;

export const UpdateDashboard = gql`
    mutation UpdateDashboard($id: Int!, $name: String!) {
        updateDashboard(id: $id, name: $name) {
            id
        }
    }
`;
export const CreateDashboard = gql`
    mutation CreateDashboard($name: String!) {
        createDashboard(name: $name) {
            id
        }
    }
`;

export const UpdatePos = gql`
    mutation UpdatePos($entryId: Int!, $pos: InputResponsiveDashboardEntryPos!) {
        updateDashboardEntry(entryId: $entryId, pos: $pos) {
            id
        }
    }
`;

export const UpdateDashboardEntry = gql`
    mutation UpdateDashboardEntry(
        $entryId: Int!
        $entryType: EntryType!
        $title: String!
        $total: Boolean!
        $stats: InputStatsSelection!
    ) {
        updateDashboardEntry(entryId: $entryId, entryType: $entryType, title: $title, total: $total, stats: $stats) {
            id
        }
    }
`;
export const AddDashboardEntry = gql`
    mutation AddDashboardEntry(
        $dashboardId: Int!
        $entryType: EntryType!
        $title: String!
        $total: Boolean!
        $stats: InputStatsSelection!
        $pos: InputResponsiveDashboardEntryPos
    ) {
        addDashboardEntry(
            dashboardId: $dashboardId
            entryType: $entryType
            title: $title
            total: $total
            stats: $stats
            pos: $pos
        ) {
            id
        }
    }
`;
export const RemoveDashboardEntry = gql`
    mutation RemoveDashboardEntry($id: Int!) {
        removeDashboardEntry(id: $id) {
            id
        }
    }
`;

export const RemoveDashboardRange = gql`
    mutation RemoveDashboardRange($rangeId: Int!) {
        removeDashboardRange(rangeId: $rangeId) {
            id
        }
    }
`;

export const UpdateDashboardRange = gql`
    mutation UpdateDashboardRange($rangeId: Int!, $range: InputNamedDateRange!) {
        updateDashboardRange(rangeId: $rangeId, range: $range) {
            name
        }
    }
`;

export const AddDashboardRange = gql`
    mutation AddDashboardRange($dashboardId: Int!, $range: InputNamedDateRange!) {
        addDashboardRange(dashboardId: $dashboardId, range: $range) {
            name
        }
    }
`;
