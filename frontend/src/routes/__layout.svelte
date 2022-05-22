<script lang="ts" context="module">
  import "../app.postcss";

	import Notifications from 'svelte-notifications';
</script>

<Notifications>
	<slot/>
</Notifications>