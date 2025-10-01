import {createApp} from 'vue'

console.log('main.js loading...');

const app = createApp({
  template: `
    <div style="padding: 2rem; text-align: center; font-family: Arial;">
      <h1>Vue App Test</h1>
      <p>If you can see this, Vue is working!</p>
      <button @click="testClick" style="padding: 1rem; background: blue; color: white; border: none; border-radius: 4px;">
        Test Button
      </button>
    </div>
  `,
  setup() {
    console.log('Vue app setup running');
    
    function testClick() {
      console.log('Button clicked!');
      alert('Vue is working!');
    }
    
    return { testClick };
  }
});

console.log('Mounting app...');
app.mount('#app');
console.log('App mounted!');