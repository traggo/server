import {Tags} from '../gql/__generated__/Tags';
import * as gqlTags from '../gql/tags';
import {QueryHookResult, useQuery} from 'react-apollo-hooks';
import {SuggestTagValue, SuggestTagValueVariables} from '../gql/__generated__/SuggestTagValue';
import {TagSelectorEntry, specialTag} from './tagSelectorEntry';

export const useSuggest = (
    tagResult: QueryHookResult<Tags, {}>,
    inputValue: string,
    usedTags: string[],
    skipValue = false
): TagSelectorEntry[] => {
    const [tagKeySomeCase, tagValue] = inputValue.split(':');
    const tagKey = tagKeySomeCase.toLowerCase();

    const exactMatch = ((tagResult.data && tagResult.data.tags) || []).find((tag) => tag.key === tagKey);

    const valueResult = useQuery<SuggestTagValue, SuggestTagValueVariables>(gqlTags.SuggestTagValue, {
        variables: {tag: tagKey, query: tagValue},
        skip: exactMatch === undefined || skipValue,
    });

    if (exactMatch && tagValue !== undefined && usedTags.indexOf(exactMatch.key) === -1 && !skipValue) {
        return suggestTagValue(exactMatch, tagValue, valueResult);
    } else {
        return suggestTag(exactMatch, tagResult, tagKey, usedTags);
    }
};

const suggestTag = (
    exactMatch: TagSelectorEntry['tag'] | undefined,
    tagResult: QueryHookResult<Tags, {}>,
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
        .sort((a, b) => a.key.localeCompare(b.key))
        .slice(0, 5)
        .map((tag) => ({tag}));
};

const suggestTagValue = (
    exactMatch: TagSelectorEntry['tag'],
    tagValue: string,
    valueResult: QueryHookResult<SuggestTagValue, SuggestTagValueVariables>
): TagSelectorEntry[] => {
    if (!valueResult.data || valueResult.data.values === null) {
        return [];
    }

    let someValues = valueResult.data.values || [];

    if (someValues.indexOf(tagValue) === -1) {
        someValues = [tagValue, ...someValues];
    }

    return someValues.map((val) => ({tag: exactMatch, value: val}));
};
