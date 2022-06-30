<script lang="ts">
  import {
    Calibrate,
    Read,
    Clear,
    GetListings,
    GetConfig,
    SetListingCount,
    SetPrice,
    SetLeague,
    Copy,
    SetName,
    SetStream,
    GetDisplayCount,
    SetDisplay,
    GetLanguages,
    SetLanguage,
    SetShortcut
  } from '../wailsjs/go/main/App';
  import type { types, main } from '../wailsjs/go/models';
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime';
  import { getNotificationsContext } from 'svelte-notifications';
  import Icon from '@iconify/svelte';
  import { t, locale, localeMapping, colorKeys } from '../lib/translations';

  const { addNotification } = getNotificationsContext();

  let reading = false;
  let calibrating = false;

  let listings: types.ParsedListing[] = [];
  let config: main.ConvertedConfig;
  let displayCount: number;
  let languages: Record<string, string>;

  let settings = false;

  let eventsInitialized = false;

  onMount(() => {
    initEvents();
    loadListings();
    loadConfig();
    GetDisplayCount().then((d) => (displayCount = d));
    GetLanguages().then((l) => (languages = l));
  });

  const initEvents = () => {
    if (eventsInitialized) {
      return;
    }

    eventsInitialized = true;

    EventsOn('listings_updated', () => {
      loadListings();
    });

    EventsOn('config_updated', () => {
      loadConfig();
    });

    EventsOn('reading_listings', () => {
      reading = true;
    });

    EventsOn('listings_read', () => {
      reading = false;
    });

    EventsOn('error', (err) => {
      handle(err);
    });

    EventsOn('warning', (warn) => {
      addNotification({
        text: warn,
        position: 'bottom-center',
        type: 'warning',
        removeAfter: 5000
      });
    });
  };

  const loadListings = () => {
    return GetListings().then((l) => {
      listings = l.sort((a, b) => {
        return a.type.localeCompare(b.type);
      });
    });
  };

  const loadConfig = () => {
    return GetConfig().then((c) => {
      config = c;
      locale.set(localeMapping[c.language]);
    });
  };

  const read = () => {
    reading = true;
    handle(Read());
  };

  const calibrate = () => {
    calibrating = true;
    handle(
      Calibrate()
        .then(() => loadConfig())
        .then(() => {
          addNotification({
            text: 'Calibrated to: ' + config.scaling,
            position: 'bottom-center',
            type: 'success',
            removeAfter: 5000
          });
        })
    ).then(() => {
      calibrating = false;
    });
  };

  const handle = (potentialError: Promise<Error> | Error) => {
    return Promise.resolve(potentialError)
      .catch((err) => err)
      .then((err) => {
        if (err) {
          console.error(err);
          addNotification({
            text: err,
            position: 'bottom-center',
            type: 'danger',
            removeAfter: 5000
          });
        }
      });
  };

  let changingShortcut = false;

  const changeShortcut = () => {
    changingShortcut = !changingShortcut;
  };

  const modKeys = ['Control', 'Alt', 'Shift', 'Meta'];

  const keyDown = (event: KeyboardEvent) => {
    if (changingShortcut) {
      if (modKeys.indexOf(event.key) < 0) {
        changingShortcut = false;
        const combination = [];
        if (event.ctrlKey) {
          combination.push('ctrl');
        }
        if (event.shiftKey) {
          combination.push('shift');
        }
        if (event.altKey || event.metaKey) {
          combination.push('alt');
        }
        combination.push(event.key);
        handle(SetShortcut(combination));
      }
    }
  };

  const changeLanguage = (lang: string) => {
    handle(SetLanguage(lang));
    locale.set(lang);
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  function colorMessages(cfg: main.ConvertedConfig, _: string) {
    const result = {};
    if (cfg) {
      Object.keys(cfg.messages).forEach((k) => {
        result[k] = colorMessage(cfg.messages[k]);
      });
    }
    return result;
  }

  $: colorMappedMessages = colorMessages(config, $locale);

  const colorMessage = (message: string): string => {
    Object.keys(colorKeys).forEach((key) => {
      const value = colorKeys[key];
      const translated = $t('colorMapping.' + key);
      translated.split(';').forEach((k) => {
        message = message.replace(
          new RegExp(`(${k}(?:$|\\s))|((?:^|\\s)${k})`, 'gi'),
          `<span style='color: ${value}; font-weight: bold'>$1$2</span>`
        );
      });
    });

    return message;
  };

  const sorted = (listings: types.ParsedListing[]): types.ParsedListing[] => {
    const copy = [...listings];
    copy.sort((a, b) => {
      const byName = a.type.localeCompare(b.type);
      if (byName != 0) {
        return byName;
      }
      return a.level - b.level;
    });
    return listings;
  };
</script>

<svelte:window on:keydown={keyDown} />

{#if config && displayCount && languages && colorMappedMessages}
  <div class="p-3 flex flex-col gap-3 w-full h-screen">
    <div class="flex flex-row flex-wrap gap-2 justify-between">
      <div>
        <h1 class="text-3xl">
          {#if !settings}
            {$t('name')}
          {:else}
            {$t('toggle.settings')}
          {/if}
        </h1>
      </div>
      {#if !settings}
        <div class="flex flex-row flex-wrap gap-2">
          <button on:click={read} class="bg-green-700 disabled:bg-green-900" disabled={reading}>
            {#if reading}
              <Icon icon="eos-icons:loading" />
            {:else}
              {$t('home.scan')}
            {/if}
          </button>
          <button on:click={() => handle(Clear())} class="bg-red-700">
            {$t('home.clear')}
          </button>
          <button on:click={() => handle(Copy())} class="bg-blue-700">
            {$t('home.copy')}
          </button>
        </div>
      {/if}
      <div>
        <button on:click={() => (settings = !settings)} class="bg-yellow-700">
          {#if !settings}
            {$t('toggle.settings')}
          {:else}
            {$t('toggle.listings')}
          {/if}
        </button>
      </div>
    </div>
    {#if !settings}
      {#if listings.length > 0}
        <div class="grid grid-flow-row gap-2 overflow-auto">
          <div class="flex flex-row gap-3">
            <span class="count-col">{$t('home.table.count')}</span>
            <span class="w-full">{$t('home.table.name')}</span>
            <span class="market-col text-center">{$t('home.table.market')}</span>
            <span class="price-col">{$t('home.table.price')}</span>
          </div>

          {#each sorted(listings) as listing}
            <div class="flex flex-row gap-3">
              <div class="flex flex-row count-col gap-2">
                <input
                  type="text"
                  value={listing.count}
                  class="w-6/12 min-w-[35px]"
                  on:blur={(event) =>
                    handle(SetListingCount(listing.type, listing.level, parseInt(event.target.value)))}
                />
                <button
                  on:click={() => handle(SetListingCount(listing.type, listing.level, listing.count + 1))}
                  class="bg-green-700 p-0 w-3/12 text-center min-w-[35px]"
                >
                  +
                </button>
                <button
                  on:click={() => handle(SetListingCount(listing.type, listing.level, listing.count - 1))}
                  class="bg-red-700 p-0 w-3/12 text-center min-w-[35px]"
                >
                  -
                </button>
              </div>
              <span style="line-height: 32px" class="w-full">
                {@html colorMappedMessages[listing.type]} ({listing.level})
              </span>
              <span style="line-height: 32px" class="market-col text-center"
                >{config.market_prices[listing.type] || '?'}</span
              >
              <input
                type="text"
                class="price-col"
                value={config.prices[listing.type]}
                on:blur={(event) => handle(SetPrice(listing.type, event.target.value))}
              />
            </div>
          {/each}
        </div>
      {:else}
        <div class="flex h-full items-center justify-center">
          {@html $t('home.scanMessage')}
        </div>
      {/if}
    {:else}
      <div class="flex flex-row gap-4 w-full">
        <div class="flex flex-col gap-4 overflow-auto w-1/3">
          <div class="flex flex-col gap-2">
            <label class="text-lg" for="language">{$t('settings.language')}</label>
            <select value={config.language} id="language" on:change={(event) => changeLanguage(event.target.value)}>
              {#each Object.keys(languages) as lang}
                <option value={lang}>{languages[lang]}</option>
              {/each}
            </select>
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-lg" for="display">{$t('settings.display')}</label>
            <select
              value={config.display}
              id="display"
              on:change={(event) => handle(SetDisplay(parseInt(event.target.value, 10)))}
            >
              {#each Array(displayCount) as _, i}
                <option value={i}>{$t('settings.display')} {i}</option>
              {/each}
            </select>
            <button
              on:click={() => calibrate()}
              class="bg-yellow-700 disabled:bg-yellow-900 text-center"
              disabled={calibrating}
            >
              {#if calibrating}
                <Icon icon="eos-icons:loading" class="inline-block" />
              {:else}
                {$t('settings.calibrate')}
              {/if}
            </button>
          </div>
        </div>
        <div class="flex flex-col gap-4 overflow-auto w-1/3">
          <div class="flex flex-row gap-2">
            <label class="text-lg" for="stream">{$t('settings.canStream')}</label>
            <input
              type="checkbox"
              checked={config.stream}
              class="w-4 h-4 self-center"
              id="stream"
              on:change={() => handle(SetStream(!config.stream))}
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-lg" for="name">{$t('settings.characterName')}</label>
            <input
              type="text"
              value={config.name}
              id="name"
              placeholder="Character Name"
              class="text-lg p-3 rounded"
              on:blur={(event) => handle(SetName(event.target.value))}
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-lg" for="league">{$t('settings.league')}</label>
            <select value={config.league} id="league" on:change={(event) => handle(SetLeague(event.target.value))}>
              <option value="std">{$t('settings.leagues.standard')}</option>
              <option value="lsc">{$t('settings.leagues.softcore')}</option>
              <option value="lhc">{$t('settings.leagues.hardcore')}</option>
            </select>
          </div>
        </div>
        <div class="flex flex-col gap-4 overflow-auto w-1/3">
          <div class="flex flex-col gap-2">
            <label class="text-lg">
              {$t('settings.shortcut')} (<code class="uppercase">{config.shortcut.join(' + ')}</code>)
            </label>
            <button on:click={() => changeShortcut()} class="bg-yellow-700 text-center">
              {#if changingShortcut}
                <Icon icon="eos-icons:loading" class="inline-block" />
              {:else}
                {$t('settings.changeShortcut')}
              {/if}
            </button>
          </div>
        </div>
      </div>
    {/if}
  </div>
{:else}
  {$t('loading')}
{/if}

<style lang="postcss">
  .count-col {
    @apply w-2/12;
    max-width: 200px;
  }

  .market-col,
  .price-col {
    @apply w-1/12;
    max-width: 150px;
  }
</style>
