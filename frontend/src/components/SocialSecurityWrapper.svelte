<script>
  // Simple wrapper around SocialSecuritySection to help debug errors
  import SocialSecuritySection from './SocialSecuritySection.svelte';
  
  const props = $props();
</script>

<div class="wrapper">
  <h3 class="mb-4">Social Security Wrapper (Debug)</h3>
  <SocialSecuritySection {...props} />
</div>