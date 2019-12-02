// @ts-ignore
import colorHash from 'color-hash';

const dark = new colorHash({saturation: 0.35, lightness: 0.35});
const darkNone = new colorHash({saturation: 0.35, lightness: 0.5});

const light = new colorHash({saturation: 0.5, lightness: 0.6});
const lightNone = new colorHash({saturation: 0.5, lightness: 0.4});

export enum ColorMode {
    Bold,
    None,
}

const mapping: Record<ColorMode, Record<'dark' | 'light', (s: string) => string>> = {
    [ColorMode.Bold]: {
        dark: (s) => dark.hex(s),
        light: (s) => light.hex(s),
    },
    [ColorMode.None]: {
        dark: (s) => darkNone.hex(s),
        light: (s) => lightNone.hex(s),
    },
};

export const calculateColor = (s: string, mode: ColorMode, theme: 'light' | 'dark') => {
    return mapping[mode][theme](s);
};
