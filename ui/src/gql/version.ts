import {gql} from '@apollo/client';
import {VersionQuery as VersionResponse} from './__generated__';

export const Version = gql`
    query Version {
        version {
            name
            commit
            buildDate
        }
    }
`;

export const VersionDefault: VersionResponse = {
    version: {__typename: 'Version', commit: 'unknown', buildDate: 'unknown', name: 'vUnknown'},
};
