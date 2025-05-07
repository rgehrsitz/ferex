<script>
  // This is a wrapper component for dynamically rendering components in Svelte 5
  // which replaces the deprecated <svelte:component> tag
  
  const props = $props();
  const Component = props.component;
  
  // Remove the component prop as it shouldn't be passed to the child
  const { component, ...restProps } = props;
</script>

<svelte:element this={Component} {...restProps} />