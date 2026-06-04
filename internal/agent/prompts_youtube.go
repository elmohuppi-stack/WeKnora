package agent

// WikiYouTubeTranscriptPrompt transforms a raw YouTube transcript into a
// well-structured, readable wiki article. The raw transcript is typically
// hard to read (speaker-dependent, filler words, no paragraph structure),
// and this prompt instructs the LLM to produce a polished article while
// preserving all factual content.
const WikiYouTubeTranscriptPrompt = `You are a wiki editor specializing in transforming video transcripts into high-quality, well-structured wiki articles.

Below is the raw transcript of a YouTube video. Your task is to rewrite it as a proper, readable wiki article.

## Raw Transcript

{{.Transcript}}

## Video Information
- **Title**: {{.Title}}
- **Channel**: {{.Channel}}
- **Duration**: {{.Duration}} minutes

## Instructions

1. **Rewrite completely**: Transform the raw spoken-word transcript into a polished, well-structured wiki article. Remove filler words, repetitions, and speech disfluencies while preserving ALL factual content, examples, and key points.

2. **Structure with headings**: Organize the content into logical sections with proper Markdown headings (## for main sections, ### for subsections). Identify natural topic transitions in the video and use them as section breaks.

3. **Format for readability**:
   - Use paragraphs instead of raw transcript lines
   - Add bullet lists or numbered lists where appropriate
   - Use **bold** for key terms and concepts
   - Include relevant code blocks (triple backticks) if technical code is discussed
   - Add block quotes (> ) for important statements or quotes from the video

4. **Internal wiki links**: Where the video mentions a well-known concept, tool, or technology, format it as a wiki link using [[concept/name|display name]] syntax. For example: "We use [[concept/retrieval-augmented-generation|RAG]] to improve search results."

5. **Content preservation**: Do NOT omit any substantive content. The article should cover everything discussed in the video, just in a much more readable form. If the video has specific examples, code snippets, or demonstrations, include them.

6. **Opening paragraph**: Start with a concise introduction (1-2 paragraphs) that explains what the video/article is about and why it matters.

7. **Closing section**: End with a "## Key Takeaways" section that lists the 3-5 most important points from the video as bullet points.

8. **Language**: Write the article in {{.Language}}.

9. **Empty transcript**: If the transcript is empty or contains no substantive content, output exactly: "No transcript content was available for this video." and nothing else.

Output ONLY the wiki article content in Markdown format. Do not include any preamble or metadata before the content.`
