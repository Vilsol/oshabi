<script lang='ts'>
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
		SetLanguage
	} from '../wailsjs/go/main/App';
	import type { types, main } from '../wailsjs/go/models';
	import { onMount } from 'svelte';
	import { EventsOn } from '../wailsjs/runtime';
	import { getNotificationsContext } from 'svelte-notifications';
	import Icon from '@iconify/svelte';

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
		GetDisplayCount().then(d => displayCount = d);
		GetLanguages().then(l => languages = l);
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
		return GetListings().then(l => {
			listings = l.sort((a, b) => {
				return a.type.localeCompare(b.type);
			});
		});
	};

	const loadConfig = () => {
		return GetConfig().then(c => config = c);
	};

	const read = () => {
		reading = true;
		handle(Read());
	};

	const calibrate = () => {
		calibrating = true;
		handle(Calibrate().then(() => loadConfig()).then(() => {
			addNotification({
				text: 'Calibrated to: ' + config.scaling,
				position: 'bottom-center',
				type: 'success',
				removeAfter: 5000
			});
		})).then(() => {
			calibrating = false;
		});
	};

	const handle = (potentialError: Promise<Error> | Error) => {
		return Promise.resolve(potentialError).catch(err => err).then(err => {
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
</script>

{#if config && displayCount && languages}
	<div class='p-3 flex flex-col gap-3 w-full h-screen'>
		<div class='flex flex-row flex-wrap gap-2 justify-between'>
			<div>
				<h1 class='text-3xl'>
					{#if !settings}
						Oshabi
					{:else}
						Settings
					{/if}
				</h1>
			</div>
			{#if !settings}
				<div class='flex flex-row flex-wrap gap-2'>
					<button on:click={read} class='bg-green-700 disabled:bg-green-900' disabled={reading}>
						{#if reading}
							<Icon icon='eos-icons:loading' />
						{:else}
							Scan
						{/if}
					</button>
					<button on:click={() => handle(Clear())} class='bg-red-700'>
						Clear
					</button>
					<button on:click={() => handle(Copy())} class='bg-blue-700'>
						Copy
					</button>
				</div>
			{/if}
			<div>
				<button on:click={() => settings = !settings} class='bg-yellow-700'>
					{#if !settings}
						Settings
					{:else}
						Listings
					{/if}
				</button>
			</div>
		</div>
		{#if !settings}
			{#if listings.length > 0}
				<div class='grid grid-flow-row gap-2 overflow-auto'>
					{#each listings as listing}
						<div class='flex flex-row gap-3'>
							<div class='flex flex-row w-2/12 gap-2'>
								<input type='text' value={listing.count} class='w-6/12'
											 on:blur={(event) => handle(SetListingCount(listing.type, listing.level, parseInt(event.target.value)))}>
								<button on:click={() => handle(SetListingCount(listing.type, listing.level, listing.count+1))}
												class='bg-green-700 py-0 w-3/12'>
									+
								</button>
								<button on:click={() => handle(SetListingCount(listing.type, listing.level, listing.count-1))}
												class='bg-red-700 py-0 w-3/12'>
									-
								</button>
							</div>
							<span style='line-height: 32px' class='w-7/12'>{config.messages[listing.type]} ({listing.level})</span>
							<span style='line-height: 32px'
										class='w-1/12 text-center'>{config.market_prices[listing.type] || '?'}</span>
							<input type='text' class='w-2/12' value={config.prices[listing.type]}
										 on:blur={(event) => handle(SetPrice(listing.type, event.target.value))}>
						</div>
					{/each}
				</div>
			{:else}
				<div class='flex h-full items-center justify-center'>
					Press&nbsp;<span class='text-green-600 font-bold'>Scan</span>&nbsp;to add listings
				</div>
			{/if}
		{:else}
			<div class='flex flex-col gap-4' style='width: fit-content'>
				<div class='flex flex-col gap-2'>
					<label class='text-lg' for='language'>Language</label>
					<select value={config.language} id='language'
									on:change={(event) => handle(SetLanguage(event.target.value))}>
						{#each Object.keys(languages) as lang}
							<option value={lang}>{lang}: {languages[lang]}</option>
						{/each}
					</select>
				</div>

				<div class='flex flex-row gap-2'>
					<label class='text-lg' for='stream'>Can Stream</label>
					<input type='checkbox' checked={config.stream} class='w-4 h-4 self-center' id='stream'
								 on:change={() => handle(SetStream(!config.stream))}>
				</div>

				<div class='flex flex-col gap-2'>
					<label class='text-lg' for='name'>Character Name</label>
					<input type='text' value={config.name} id='name' placeholder='Character Name' class='text-lg p-3 rounded'
								 on:blur={(event) => handle(SetName(event.target.value))}>
				</div>

				<div class='flex flex-col gap-2'>
					<label class='text-lg' for='league'>League</label>
					<select value={config.league} id='league'
									on:change={(event) => handle(SetLeague(event.target.value))}>
						<option value='std'>Standard</option>
						<option value='lsc'>League Softcore</option>
						<option value='lhc'>League Hardcore</option>
					</select>
				</div>

				<div class='flex flex-col gap-2'>
					<label class='text-lg' for='display'>Display</label>
					<select value={config.display} id='display'
									on:change={(event) => handle(SetDisplay(parseInt(event.target.value, 10)))}>
						{#each Array(displayCount) as _, i}
							<option value={i}>Display {i}</option>
						{/each}
					</select>
					<button on:click={() => calibrate()} class='bg-yellow-700 disabled:bg-yellow-900' disabled={calibrating}>
						{#if calibrating}
							<Icon icon='eos-icons:loading' />
						{:else}
							Calibrate
						{/if}
					</button>
				</div>
			</div>
		{/if}
	</div>
{:else}
	Loading...
{/if}
