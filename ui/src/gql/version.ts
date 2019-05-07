import {gql} from 'apollo-boost';
import {Version as VersionResponse} from './__generated__/Version';

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
