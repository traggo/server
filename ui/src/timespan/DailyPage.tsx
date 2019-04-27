import * as React from 'react';
import {Tracker} from './Tracker';
import {ActiveTrackers} from './ActiveTrackers';
import {DoneTrackers} from './DoneTrackers';

export const DailyPage = () => {
    return (
        <div style={{margin: '1px auto', maxWidth: 1000}}>
            <Tracker />
            <ActiveTrackers />
            <DoneTrackers />
        </div>
    );
};
