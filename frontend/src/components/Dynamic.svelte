<script lang="ts">
  // Idiomatic Svelte 5 runes: dynamic component wrapper
  const props = $props();
  const Component = props.component;
  // Remove the component prop as it shouldn't be passed to the child
  const { component, ...restProps } = props;
  // Debug: log what is being passed to the dynamic child
  try {
    console.log('Dynamic.svelte forwarding to child:', JSON.stringify({ Component: Component?.name, ...restProps }, null, 2));
  } catch (e) {
    console.log('Dynamic.svelte forwarding to child (raw):', { Component, restProps });
  }
</script>

{#if typeof Component === 'string'}
  <svelte:element this={Component} {...restProps} />
{:else}
  <Component {...restProps} />
{/if}