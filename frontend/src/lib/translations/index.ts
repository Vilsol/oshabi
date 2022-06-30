import i18n from '@sveltekit-i18n/base';
import parser from '@sveltekit-i18n/parser-default';
import type { Config } from '@sveltekit-i18n/parser-default';

import languages from './languages.json';

import en from './en.json';
import pt from './pt.json';
import ru from './ru.json';
import th from './th.json';
import fr from './fr.json';
import de from './de.json';
import es from './es.json';
import zh from './zh.json';
import ko from './ko.json';
import ja from './ja.json';

const config: Config = {
  initLocale: 'en',
  fallbackLocale: 'en',

  parser: parser(),

  translations: {
    en: {
      ...languages,
      ...en
    },
    pt: {
      ...languages,
      ...pt
    },
    ru: {
      ...languages,
      ...ru
    },
    th: {
      ...languages,
      ...th
    },
    fr: {
      ...languages,
      ...fr
    },
    de: {
      ...languages,
      ...de
    },
    es: {
      ...languages,
      ...es
    },
    zh: {
      ...languages,
      ...zh
    },
    ko: {
      ...languages,
      ...ko
    },
    ja: {
      ...languages,
      ...ja
    }
  }
};

export const { t, locale } = new i18n(config);

export const localeMapping = {
  eng: 'en',
  por: 'pt',
  rus: 'ru',
  tha: 'th',
  fra: 'fr',
  deu: 'de',
  spa: 'es',
  chi_sim: 'zh',
  kor: 'ko',
  jpn: 'ja'
};

export const colorKeys = {
  red: '#c8676e',
  blue: '#a2cffb',
  green: '#86bda3',
  physical: '#c79d93',
  caster: '#b3f8fe',
  fire: '#ff9a77',
  cold: '#93d8ff',
  lightning: '#f8cb76',
  attack: '#da814d',
  life: '#c96e6e',
  speed: '#cfeea5',
  defence: '#a88f67',
  chaos: '#d8a7d3',
  unique: '#af6025',
  critical: '#b2a7d6',
  magicItem: '#8888ff'
};
