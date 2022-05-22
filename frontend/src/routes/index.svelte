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
		SetDisplay
	} from '../wailsjs/go/main/App';
	import type { types, main } from '../wailsjs/go/models';
	import { onMount } from 'svelte';
	import { EventsOn } from '../wailsjs/runtime';
	import { getNotificationsContext } from 'svelte-notifications';

	const { addNotification } = getNotificationsContext();

	let reading = false;

	let listings: types.ParsedListing[] = [];
	let config: main.ConvertedConfig;
	let displayCount: number;

	let eventsInitialized = false;

	onMount(() => {
		initEvents();
		loadListings();
		loadConfig();
		GetDisplayCount().then(d => displayCount = d);
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
		GetConfig().then(c => config = c);
	};

	const read = () => {
		reading = true;
		handle(Read());
	};

	const handle = (potentialError: Promise<Error> | Error) => {
		Promise.resolve(potentialError).catch(err => err).then(err => {
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

{#if config && displayCount}
	<div class='p-3 flex flex-col gap-3 w-full h-screen'>
		<div class='flex flex-row flex-wrap justify-between'>
			<div class='flex flex-row gap-2'>
				<button on:click={read} class='bg-green-700 disabled:bg-green-900' disabled={reading}>
					Scan
				</button>
				<button on:click={() => handle(Clear())} class='bg-red-700'>
					Clear
				</button>
				<button on:click={() => handle(Copy())} class='bg-blue-700'>
					Copy
				</button>
			</div>
			<div class='flex flex-row gap-2'>
				<div class='flex flex-col text-center'>
					<label for='stream'>Stream</label>
					<input type='checkbox' value={config.stream} class='w-4 h-4 m-auto' id='stream'
								 on:change={() => handle(SetStream(!config.stream))}>
				</div>
				<input type='text' value={config.name} placeholder='Name' class='text-lg p-3 rounded'
							 on:blur={(event) => handle(SetName(event.target.value))}>
				<select value={config.league} on:change={(event) => handle(SetLeague(event.target.value))}>
					<option value='std'>Standard</option>
					<option value='lsc'>League Softcore</option>
					<option value='lhc'>League Hardcore</option>
				</select>
				<select value={config.display} on:change={(event) => handle(SetDisplay(parseInt(event.target.value, 10)))}>
					{#each Array(displayCount) as _, i}
						<option value={i}>Display {i}</option>
					{/each}
				</select>
				<button on:click={() => handle(Calibrate())} class='bg-yellow-700'>
					Calibrate
				</button>
			</div>
		</div>
		{#if listings.length > 0}
			<div class='grid grid-flow-row gap-2'>
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
	</div>
{:else}
	Loading...
{/if}
