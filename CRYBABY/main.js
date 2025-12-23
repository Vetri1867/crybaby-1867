
(function(){
      // Lightweight app state using localStorage to persist (demo cloud-like)
      const store = {
        get(key, fallback){
          try{
            const raw = localStorage.getItem(key);
            return raw ? JSON.parse(raw) : fallback;
          }catch(e){return fallback}
        },
        set(key, val){ localStorage.setItem(key, JSON.stringify(val)); },
        remove(key){ localStorage.removeItem(key); }
      };

      // ... (rest of your existing code)

      async function getAIResponse(promptText, mode){
        // If useRealApi checked and key provided, attempt to call OpenAI's chat completions endpoint (mock safety: small)
        if(useRealApi.checked && apiKeyInput.value.trim()){
          try{
            const apiKey = apiKeyInput.value.trim();
            const body = {
              model: "gpt-4o-mini", // placeholder, users may change
              messages: [{role:'user', content: `Explain in simple child-friendly language: ${promptText}. Use emojis and encouraging tone.`}],
              max_tokens: 400
            };
            const res = await fetch('https://api.openai.com/v1/chat/completions', {
              method:'POST', headers:{ 'Content-Type':'application/json', 'Authorization':'Bearer '+apiKey }, body: JSON.stringify(body)
            });
            const js = await res.json();
            const content = js.choices && js.choices[0] && (js.choices[0].message?.content || js.choices[0].text) || 'Sorry, API did not return a response.';
            return content;
          }catch(err){
            console.error(err);
            return "I couldn't reach the real API. I'll help with my built-in tutor instead 😊";
          }
        }
        
        // Use the Go backend API for the AI tutor
        try {
          const response = await fetch('/api/tutor', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ prompt: promptText }),
          });

          if (!response.ok) {
            throw new Error('API request failed');
          }

          const data = await response.json();
          // Extract the text from the response. This may need adjustment based on the actual Gemini API response structure.
          return data.candidates[0].content.parts[0].text;
        } catch (error) {
          console.error('Error calling Go backend:', error);
          return "I'm having a little trouble connecting to my brain right now. Please try again in a moment! 🧸";
        }
      }

      // ... (rest of your existing code)
    })();
