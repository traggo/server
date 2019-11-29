import {Tags, Tags_tags} from '../gql/__generated__/Tags';
import {TagDefinitionType} from '../gql/__generated__/globalTypes';
import {QueryResult} from 'react-apollo';

export interface TagInputError {
    error: string;
    value: string;
}

export interface TagSelectorEntry {
    tag: Omit<Tags_tags, 'usages'> & {create?: boolean; alreadyUsed?: boolean};
    value?: string;
}

export interface InputTag {
    key: string;
    stringValue: string | null;
}

export const toInputTags = (entries: TagSelectorEntry[]): InputTag[] => {
    return entries.map((entry) => ({key: entry.tag.key, stringValue: entry.value || null}));
};

export const toTagSelectorEntry = (tags: Array<TagSelectorEntry['tag']>, entries: InputTag[]): TagSelectorEntry[] => {
    return entries.map(
        (timerTag): TagSelectorEntry => {
            const definition = tags.find((tag) => tag.key === timerTag.key) || specialTag(timerTag.key, 'new');
            return {
                tag: definition,
                value: timerTag.stringValue || undefined,
            };
        }
    );
};

export const specialTag = (name: string, state: 'used' | 'new'): TagSelectorEntry['tag'] & {usages: 0} => {
    return {
        key: name,
        __typename: 'TagDefinition',
        color: 'gray',
        type: TagDefinitionType.novalue,
        create: state === 'new',
        alreadyUsed: state === 'used',
        usages: 0,
    };
};

const tryAdd = (tagsResult: QueryResult<Tags, {}>, entry: string): TagSelectorEntry | TagInputError => {
    const [keySomeCase, value, ...other] = entry.split(':');
    const key = keySomeCase.toLowerCase();

    if (other.length || !tagsResult.data || !tagsResult.data.tags) {
        return {error: `${entry} has too many colons`, value: entry};
    }

    const foundTag = tagsResult.data.tags.find((tag) => tag.key === key);

    if (foundTag === undefined) {
        return {error: `'${key}' does not exist`, value: entry};
    }

    if (foundTag.type === TagDefinitionType.singlevalue && !value) {
        return {error: `'${key}' requires a value`, value: entry};
    }

    return {tag: foundTag, value};
};
interface EntriesAndErrors {
    errors: TagInputError[];
    entries: TagSelectorEntry[];
    usedTags: string[];
}

const groupAndCheckExistence = (a: EntriesAndErrors, entry: TagInputError | TagSelectorEntry) => {
    if ('tag' in entry) {
        if (a.usedTags.indexOf(entry.tag.key) === -1) {
            a.entries = [...a.entries, entry];
            a.usedTags = [...a.usedTags, entry.tag.key];
        } else {
            a.errors = [...a.errors, {value: label(entry), error: `'${entry.tag.key}' is already defined`}];
        }
    } else {
        a.errors = [...a.errors, entry];
    }
    return a;
};

export const addValues = (
    newValue: string,
    tagsResult: QueryResult<Tags, {}>,
    selectedEntries: TagSelectorEntry[]
): EntriesAndErrors => {
    return newValue
        .split(/\s+/)
        .filter((entry) => entry)
        .map((entry) => tryAdd(tagsResult, entry))
        .reduce(groupAndCheckExistence, {errors: [], entries: [], usedTags: selectedEntries.map((entry) => entry.tag.key)});
};

export const itemLabel = (tag: TagSelectorEntry) => {
    if (tag.tag.create) {
        return `Create tag '${tag.tag.key}'`;
    }
    if (tag.tag.alreadyUsed) {
        return `Tag '${tag.tag.key}' is already defined`;
    }
    return label(tag);
};

export const label = (tag: TagSelectorEntry) => {
    const suffix = tag.tag.type === TagDefinitionType.novalue ? '' : ':' + (tag.value || '');
    return `${tag.tag.key}${suffix}`;
};
