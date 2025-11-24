import {Tags} from '../gql/__generated__/Tags';
import * as gqlTags from '../gql/tags';
import {useQuery} from '@apollo/react-hooks';
import {SuggestTagValue, SuggestTagValueVariables} from '../gql/__generated__/SuggestTagValue';
import {TagSelectorEntry, specialTag} from './tagSelectorEntry';
import {QueryResult} from 'react-apollo';

export const useSuggest = (
    tagResult: QueryResult<Tags, {}>,
    inputValue: string,
    usedTags: string[],
    skipValue = false,
    includeInputValueOnNoMatch = true
): TagSelectorEntry[] => {
    const [tagKeySomeCase, tagValue] = inputValue.split(':');
    const tagKey = tagKeySomeCase.toLowerCase();

    const exactMatch = ((tagResult.data && tagResult.data.tags) || []).find((tag) => tag.key === tagKey);

    const valueResult = useQuery<SuggestTagValue, SuggestTagValueVariables>(gqlTags.SuggestTagValue, {
        variables: {tag: tagKey, query: tagValue},
        skip: exactMatch === undefined || skipValue,
    });

    if (exactMatch && tagValue !== undefined && usedTags.indexOf(exactMatch.key) === -1 && !skipValue) {
        return suggestTagValue(exactMatch, tagValue, valueResult, includeInputValueOnNoMatch);
    } else {
        return suggestTag(exactMatch, tagResult, tagKey, usedTags);
    }
};

const suggestTag = (
    exactMatch: TagSelectorEntry['tag'] | undefined,
    tagResult: QueryResult<Tags, {}>,
    tagKey: string,
    usedTags: string[]
) => {
    if (!tagResult.data || tagResult.data.tags === null) {
        return [];
    }

    let availableTags = (tagResult.data.tags || [])
        .filter((tag) => usedTags.indexOf(tag.key) === -1)
        .filter((tag) => tag.key.indexOf(tagKey) === 0);

    if (tagKey && !exactMatch) {
        availableTags = [specialTag(tagKey, 'new'), ...availableTags];
    }

    if (usedTags.indexOf(tagKey) !== -1) {
        availableTags = [specialTag(tagKey, 'used'), ...availableTags];
    }

    return availableTags
        .sort((a, b) => b.usages - a.usages)
        .slice(0, 5)
        .map((tag) => ({tag, value: ''}));
};

const suggestTagValue = (
    exactMatch: TagSelectorEntry['tag'],
    tagValue: string,
    valueResult: QueryResult<SuggestTagValue, SuggestTagValueVariables>,
    includeInputValueOnNoMatch: boolean
): TagSelectorEntry[] => {
    let someValues = (valueResult.data && valueResult.data.values) || [];

    if (includeInputValueOnNoMatch && someValues.indexOf(tagValue) === -1) {
        someValues = [tagValue, ...someValues];
    }

    return someValues.map((val) => ({tag: exactMatch, value: val}));
};
