import {gql} from 'apollo-boost';

export const Devices = gql`
    query Devices {
        devices {
            id
            name
            expiresAt
            createdAt
            activeAt
        }
        currentDevice {
            id
        }
    }
`;
export const RemoveDevice = gql`
    mutation RemoveDevice($id: Int!) {
        removeDevice(id: $id) {
            id
        }
    }
`;
export const UpdateDevice = gql`
    mutation UpdateDevice($id: Int!, $name: String!) {
        updateDevice(id: $id, name: $name) {
            id
        }
    }
`;
