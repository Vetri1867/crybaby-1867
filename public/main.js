
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

          const result = await response.json();
          if (result.error) {
            throw new Error(result.error);
          }
          // Extract the text from the response.
          return result.data.candidates[0].content.parts[0].text;
        } catch (error) {
          console.error('Error calling Go backend:', error);
          return "I'm having a little trouble connecting to my brain right now. Please try again in a moment! 🧸";
        }
      }
      
      async function searchYouTube(query) {
        try {
          const response = await fetch(`/api/youtube?q=${encodeURIComponent(query)}`);
          const result = await response.json();
          if (result.error) {
            throw new Error(result.error);
          }
          displayVideos(result.data.items);
        } catch (error) {
          console.error('Error searching YouTube:', error);
        }
      }

      function displayVideos(videos) {
        const videoResults = document.getElementById('videoResults');
        videoResults.innerHTML = '';
        videos.forEach(video => {
          const videoElement = document.createElement('div');
          videoElement.innerHTML = `
            <a href="https://www.youtube.com/watch?v=${video.id.videoId}" target="_blank">
              <img src="${video.snippet.thumbnails.default.url}" alt="${video.snippet.title}">
              <p>${video.snippet.title}</p>
            </a>
          `;
          videoResults.appendChild(videoElement);
        });
      }

      async function createCalendarEvent() {
        const summary = document.getElementById('eventSummary').value;
        const description = document.getElementById('eventDescription').value;
        const start = document.getElementById('eventStart').value;
        const end = document.getElementById('eventEnd').value;

        if (!summary || !start || !end) {
            alert('Please fill out all fields to create an event.');
            return;
        }

        try {
            const response = await fetch('/api/calendar/event', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                  summary, 
                  description, 
                  start: new Date(start).toISOString(), 
                  end: new Date(end).toISOString() 
                }),
            });

            const result = await response.json();

            if (result.error) {
                if (response.status === 401) {
                    alert(result.error);
                    window.location.href = '/login/google';
                } else {
                    throw new Error(result.error);
                }
            } else {
                alert('Event created successfully!');
            }
        } catch (error) {
            console.error('Error creating calendar event:', error);
            alert('Failed to create event. See console for details.');
        }
    }


      document.addEventListener('DOMContentLoaded', () => {
        const youtubeSearchBtn = document.getElementById('youtubeSearchBtn');
        const youtubeSearchQuery = document.getElementById('youtubeSearchQuery');
        const loginGoogleBtnStudent = document.getElementById('loginGoogleBtnStudent');
        const loginGoogleBtnParent = document.getElementById('loginGoogleBtnParent');
        const createEventBtn = document.getElementById('createEventBtn');
        const createEventForm = document.getElementById('createEventForm');

         // Check for Google Auth success flag
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.get('google_auth_success')) {
            store.set('google_auth_success', true);
            // Clean the URL
            window.history.replaceState({}, document.title, "/");
        }

        // Show/hide calendar form based on login state
        if (store.get('google_auth_success')) {
            if (createEventForm) {
                createEventForm.classList.remove('hidden');
            }
            if (loginGoogleBtnParent) {
                loginGoogleBtnParent.classList.add('hidden');
            }
        }

        if (youtubeSearchBtn) {
            youtubeSearchBtn.addEventListener('click', () => {
                const query = youtubeSearchQuery.value;
                if (query) {
                    searchYouTube(query);
                }
            });
        }

        if (loginGoogleBtnStudent) {
            loginGoogleBtnStudent.addEventListener('click', () => {
                window.location.href = '/login/google';
            });
        }

        if (loginGoogleBtnParent) {
            loginGoogleBtnParent.addEventListener('click', () => {
                window.location.href = '/login/google';
            });
        }

        if (createEventBtn) {
            createEventBtn.addEventListener('click', createCalendarEvent);
        }
      });

      // ... (rest of your existing code)
    })();
