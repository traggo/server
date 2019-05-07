import * as React from 'react';
import {Tracker} from './Tracker';
import {ActiveTrackers} from './ActiveTrackers';
import {DoneTrackers} from './DoneTrackers';
import {TagSelectorEntry} from '../tag/tagSelectorEntry';

export const DailyPage = () => {
    const [selectedEntries, setSelectedEntries] = React.useState<TagSelectorEntry[]>([]);
    return (
        <div style={{margin: '1px auto', maxWidth: 1000}}>
            <Tracker selectedEntries={selectedEntries} onSelectedEntriesChanged={setSelectedEntries} />
            <ActiveTrackers />
            <DoneTrackers
                addTagsToTracker={
                    selectedEntries.length === 0 ? (entries) => setSelectedEntries(selectedEntries.concat(entries)) : undefined
                }
            />
        </div>
    );
};
