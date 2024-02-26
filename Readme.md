# Combine two concepts into a new concept

_I wanted to see if I could build / reverse engineer Neil.fun's Infinite Crafting game (just the API portion)._

My Golang API will accept two concepts, and return the "logical combination" of the two things, plus an emoji that should make sense.

__Examples:__

- 'Water' + 'Fire' = "Steam 💨"
- 'Wind' + 'Water' = "Wave 🌊"
- 'Fire' + 'Wood' = "Smoke 💨"
- 'Wave' + 'Wind' = "Storm ⛈️"
- 'Alien' + 'Ship' = "UFO 🛸"
- 'Lightning' + 'God' = "Zeus 🤴"
- 'Assassin' + 'Whale' = "Orca 🐋"
- 'World' + 'Orc' = "Warcraft 🎮"
- 'Fantasy' + 'Bilbo' = "Hobbit 🧝"
- 'Hacker' + 'Hoodie' = "Hacktivist 💻"
- 'Mouse' + 'Food' = "Cheese 🧀"

![Screenshot](./misc/readme_images/screenshot.png)

## Env

You need a `.env` file which has:

```env
OPENAI_API_KEY=YourAPIKeyFromOpenAI
# Example: gpt-4-1106-preview
OPENAI_MODEL=YourOpenAIModel
```

- [OpenAI API Key](https://platform.openai.com/api-keys)

## Run the app

```bash
# Clone the repo, and cd into it
# ...

# Create .env file and add the API key
echo OPENAI_API_KEY=MyAPIKey > .env
echo OPENAI_MODEL=gpt-4-1106-preview >> .env

# Run the code
go run ./cmd/main.go
```

__For building the app into an exe, I used `make` with a `makefile`.__

## My Solution

- The user submits a request to my API containing two concepts
  - GET localhost:8080?one=fire&two=water
- First, I check a Redis database to see if I have seen this "Fire" + "Water" combination before:
  - If I find an existing match within Redis:
    - Then I can just return that answer to the user. This has the added benefit of ensuring all users have identical results. Plus it's extremely fast.
  - If I __DO NOT__ find an existing match within Redis:
    - Then I need to make a request to OpenAI (ChatGPT)...
    - I craft a chat request to an OpenAI model (either gpt-3.5-turbo which I have trained, or gpt-4-turbo)
    - I get the response back from ChatGPT, cache it in Redis, then respond to the user informing them that they are the first person to make that combination. _(I know they're the first, because I didn't have it in Redis already)_

__Regarding OpenAI:__

I had to [choose a model](https://openai.com/pricing) to use (to make requests against):

- `gpt-3.5-turbo-0125` is pretty dumb, but it's very cheap to use
  - $0.0005 / 1K tokens, $0.0015 / 1K tokens
- `gpt-4-1106-preview` (GPT-4 Turbo) was very good, but it was considerably more expensive
  - $0.01 / 1K tokens, $0.03 / 1K tokens
  - _"GPT-4 Turbo" is better and cheaper than just regular GPT-4_

I figured out that I can ["Fine Tune"](https://platform.openai.com/finetune) gpt-3.5-turbo with my [own training data (see link to file)](./misc/open_ai_training_data/training.jsonl). Training gpt-3.5-turbo in this way is pretty costly at first, but cheaper in the long run if you assume that users end up creating 100,000 new combinations. This is because the base model is gpt-3.5 which is much cheaper than gpt-4. If gpt-4 were cheaper, it'd be the better choice for sure.

The most difficult part of this project was figuring out how to prompt ChatGPT to give me back the right answer in a simple format. I wanted "Steam 💨". _(One word, space, one emoji)_ It kept wanting to give me 2 emojis rather than 1, or it would simply join the words together, or it'd give me a long sentence. If you're thinking, "just give it a really long detailed prompt with examples", that's what I thought too, except you are charged based on how many "tokens" you use. So, if your request takes 250 tokens total, then gpt-4 would charge you 1 penny for every 5 requests. Meaning, some jerk with a permutation script could rack up $1,000 for me _(and I saw that people were doing that to Neil.fun)_.

## Notes about the API

- I'm using Golang v1.22. The server is just using the builtin `net/http`. I didn't see a need for a framework.