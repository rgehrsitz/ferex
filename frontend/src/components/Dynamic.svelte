<script lang="ts">
  // Idiomatic Svelte 5 runes: dynamic component wrapper
  const props = $props();
  const Component = props.component;
  // Remove the component prop as it shouldn't be passed to the child
  const { component, ...restProps } = props;
  // Logging removed to prevent visible output
</script>

{#if typeof Component === 'string'}
  <svelte:element this={Component} {...restProps} />
{:else}
  <Component {...restProps} />
{/if}