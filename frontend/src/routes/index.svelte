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
		SetStream
	} from '../wailsjs/go/main/App';
	import type { types, main } from '../wailsjs/go/models';
	import { onMount } from 'svelte';
	import { EventsOn } from '../wailsjs/runtime';

	let reading = false;

	let listings: types.ParsedListing[] = [];
	let config: main.ConvertedConfig;

	onMount(() => {
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

		loadListings();
		loadConfig();
	});

	const loadListings = () => {
		return GetListings().then(l => {
			listings = l.sort((a, b) => {
				return a.type.localeCompare(b.type);
			});
		});
	};

	const loadConfig = () => {
		GetConfig().then(c => config = c).then(() => console.log(config));
	};

	const read = async () => {
		reading = true;
		await Read();
	};
</script>

{#if config}
	<div class='p-3 flex flex-col gap-3 w-full h-screen'>
		<div class='grid grid-flow-col'>
			<div class='flex flex-row gap-2'>
				<button on:click={read} class='bg-green-700 disabled:bg-green-900' disabled={reading}>
					Scan
				</button>
				<button on:click={Clear} class='bg-red-700'>
					Clear
				</button>
				<button on:click={Copy} class='bg-blue-700'>
					Copy
				</button>
			</div>
			<div class='justify-self-end flex flex-row gap-2'>
				<div class='flex flex-col text-center'>
					<label for='stream'>Stream</label>
					<input type='checkbox' value={config.stream} class='w-4 h-4 m-auto' id='stream'
								 on:change={() => SetStream(!config.stream)}>
				</div>
				<input type='text' value={config.name} placeholder='Name' class='text-lg p-3 rounded'
							 on:blur={(event) => SetName(event.target.value)}>
				<select value={config.league} on:change={(event) => SetLeague(event.target.value)}>
					<option value='std'>Standard</option>
					<option value='lsc'>League Softcore</option>
					<option value='lhc'>League Hardcore</option>
				</select>
				<button on:click={Calibrate} class='bg-yellow-700'>
					Calibrate
				</button>
			</div>
		</div>
		<div class='grid grid-flow-row gap-2'>
			{#each listings as listing}
				<div class='flex flex-row gap-3'>
					<div class='flex flex-row w-2/12 gap-2'>
						<input type='text' value={listing.count} class='w-6/12'
									 on:blur={(event) => SetListingCount(listing.type, listing.level, parseInt(event.target.value))}>
						<button on:click={() => SetListingCount(listing.type, listing.level, listing.count+1)}
										class='bg-green-700 py-0 w-3/12'>
							+
						</button>
						<button on:click={() => SetListingCount(listing.type, listing.level, listing.count-1)}
										class='bg-red-700 py-0 w-3/12'>
							-
						</button>
					</div>
					<span style='line-height: 32px' class='w-7/12'>{config.messages[listing.type]} ({listing.level})</span>
					<span style='line-height: 32px' class='w-1/12 text-center'>{config.market_prices[listing.type] || '?'}</span>
					<input type='text' class='w-2/12' value={config.prices[listing.type]}
								 on:blur={(event) => SetPrice(listing.type, event.target.value)}>
				</div>
			{/each}
		</div>
	</div>
{:else}
	Loading...
{/if}
