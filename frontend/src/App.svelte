<script lang="ts">
  import { writable } from 'svelte/store'
  import { SvelteFlow, Background, Controls } from '@xyflow/svelte'
  import { onMount } from "svelte"
  import '@xyflow/svelte/dist/style.css'
  import { GetFiles } from "../wailsjs/go/main/App"
  import { main } from "../wailsjs/go/models"

  const nodes = writable([])
  const edges = writable([])

  onMount(() => {
    const fetchData = async () => {
      const revs: main.RevFileResponse = await GetFiles()
      console.log(revs)
      nodes.set(revs.nodes)
      edges.set(revs.edges)
    }
    fetchData()
  })

</script>

<main>
  <button on:click={() => console.log($nodes)}>
    HELLO
  </button>
  <SvelteFlow {nodes} {edges} fitView>
    <Background bgColor="rgba(126,159,219,0.5)" patternColor="white" />
    <Controls />
  </SvelteFlow>
</main>

<style>
  main {
    height: 100vh;
  }
</style>
