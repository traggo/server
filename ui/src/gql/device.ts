import {gql} from 'apollo-boost';

export const Devices = gql`
    query Devices {
        devices {
            id
            name
            type
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
    mutation UpdateDevice($id: Int!, $name: String!, $deviceType: DeviceType!) {
        updateDevice(id: $id, name: $name, type: $deviceType) {
            id
        }
    }
`;

export const CreateDevice = gql`
    mutation CreateDevice($name: String!, $deviceType: DeviceType!) {
        device: createDevice(name: $name, type: $deviceType) {
            token
        }
    }
`;
