# Agent Systems Research Notes

This file is a practical, research-backed field guide for building agentic systems.

Use it for:
- prompt design
- orchestration decisions
- tool-use policy
- automation loop design
- eval and reliability practices

The goal is not to collect every paper. The goal is to keep the highest-signal findings that actually change how an agent or scaffold should be built.

## Fast Takeaways

1. Start with the simplest scaffold that can pass evals. Do not default to multi-agent.
2. Use tools when the task depends on current facts, exact details, or environment feedback.
3. Grade outcomes, artifacts, and grounded evidence, not exact tool-call traces.
4. Separate cheap mechanical work from expensive reasoning.
5. Use reflection/revision only when it improves measured performance more than it hurts latency/cost.
6. Keep prompts short, constraint-like, and verification-oriented. Avoid persona-heavy prompt sludge.
7. Read transcripts. If metrics and transcripts disagree, the harness or grader may be wrong.
8. Heterogeneous systems beat piles of homogeneous agents when the roles are genuinely different.
9. External feedback beats self-confidence. Tests, search results, compiler output, and graders matter.
10. Narrow loops outperform vague autonomy. Small mutable surface, fixed metric, bounded retries.

## What To Copy Into Systems

### Prompting

- State objective, constraints, and success criteria explicitly.
- Preserve exact terms from the user or evidence; do not rename concrete entities.
- Prefer a short rule like "if not verified, say so" over long keyword lists and examples.
- Give the model an honest escape hatch: `unknown`, `need evidence`, or `search more`.
- Use prompt tricks as baselines first, not as substitutes for retrieval, tests, or evals.

### Orchestration

- Keep the default path single-agent or workflow-based.
- Add planners, reviewers, or specialist agents only when evals show clear gains.
- Prefer bounded loops: one plan, one act phase, one verifier, one retry budget.
- Use different models or prompts only when they contribute distinct evidence or skills.
- Treat multi-agent diversity as a tool, not a religion.

### Tooling

- Favor tools that return verifiable feedback:
  - tests
  - compiler errors
  - search results
  - fetched pages
  - graders
- Keep traces and artifacts.
- Persist compact research notes when follow-up questions are common.
- If the task is stale/exact/source-sensitive, lookup beats memory.

### Automation

- Fix a metric before running an autonomous loop.
- Keep the mutable surface small.
- Auto-commit only after checks pass.
- Separate "experiment failed" from "checks failed" from "metric regressed".
- Prefer narrow optimization targets over grand autonomous platform behavior.

### Evaluation

- Build evals from real failures and real manual checks.
- Balance both sides of decision boundaries:
  - should search
  - should not search
- Isolate trials. No shared repo state, hidden cache, or leaked history.
- Use deterministic graders where possible.
- Use LLM graders with clear rubrics and human calibration when needed.
- Track both quality and consistency.

## Core Sources

## Prompting And Reasoning

### 1. Large Language Models are Zero-Shot Reasoners (Kojima et al., 2022)
Source:
- https://arxiv.org/abs/2205.11916

Why it matters:
- A very small reasoning cue can unlock much better performance than a plain direct answer.

Key takeaway:
- Before inventing a complicated prompt chain, test a minimal reasoning baseline.

Implication for agents:
- Use simple reasoning scaffolds as the baseline to beat.
- If a complex workflow does not outperform a minimal prompt + tool loop, the workflow is probably not worth it.

### 2. Self-Consistency Improves Chain of Thought Reasoning in Language Models (Wang et al., 2022)
Source:
- https://arxiv.org/abs/2203.11171

Why it matters:
- Sampling multiple reasoning paths and aggregating them can improve correctness on hard reasoning tasks.

Key takeaway:
- Best-of-N / vote / pass@k style decoding is often more useful than one brittle "perfect prompt".

Implication for agents:
- Use this selectively for high-value reasoning or planning steps.
- Do not apply it blindly to every turn; it is a latency and cost tradeoff.

### 3. ReAct: Synergizing Reasoning and Acting in Language Models (Yao et al., 2022)
Source:
- https://arxiv.org/abs/2210.03629

Why it matters:
- ReAct formalized the now-standard pattern of interleaving reasoning with external actions.

Key takeaway:
- Reasoning is better when it can touch the world.

Implication for agents:
- For factual, interactive, or environment-dependent tasks, combine thinking with acting instead of pushing all work into one monologue.
- This is a strong default for search, repo work, shell use, and structured tool loops.

### 4. Tree of Thoughts / Search-Style Deliberation
Sources:
- Tree of Thoughts: https://arxiv.org/abs/2305.10601
- Language Agent Tree Search (LATS): https://arxiv.org/abs/2310.04406

Why it matters:
- Search over candidate reasoning paths can improve hard problems when a single left-to-right pass is too brittle.

Key takeaway:
- Deliberate search can help, but it is expensive. Use it for genuinely hard branches, not for routine chat.

Implication for agents:
- Keep search/planning loops bounded.
- Reach for tree search only when the task is hard enough and the eval gain justifies the extra cost.

## Tool Use And Self-Correction

### 5. Toolformer: Language Models Can Teach Themselves to Use Tools (Schick et al., 2023)
Source:
- https://arxiv.org/abs/2302.04761

Why it matters:
- Tool use is not just a static external heuristic; it can be learned and integrated into the model's behavior.

Key takeaway:
- A good agent should decide when to call tools from the information need, not from crude keyword triggers.

Implication for agents:
- Prefer model-directed tool decisions over brittle word lists.
- Keep a simple fallback policy, but do not let the fallback dominate the product behavior.

### 6. CRITIC: Large Language Models Can Self-Correct with Tool-Interactive Critiquing (Gou et al., 2024)
Source:
- https://arxiv.org/abs/2305.11738

Why it matters:
- Self-correction becomes much stronger when the critique is grounded in external tools instead of pure self-reflection.

Key takeaway:
- Verification works better with evidence than with vibes.

Implication for agents:
- When possible, critique drafts against search results, tests, or environment state.
- A grounded revision pass is usually higher value than another creative generation pass.

### 7. Self-Refine: Iterative Refinement with Self-Feedback (Madaan et al., 2023)
Source:
- https://arxiv.org/abs/2303.17651

Why it matters:
- Even without external tools, generate -> critique -> revise can improve outputs.

Key takeaway:
- Revision is a useful primitive, but should be bounded and measured.

Implication for agents:
- Keep self-refine loops short.
- Prefer one clear revision pass over open-ended introspection.

### 8. Reflexion: Language Agents with Verbal Reinforcement Learning (Shinn et al., 2023)
Source:
- https://arxiv.org/abs/2303.11366

Why it matters:
- Reflection across attempts can improve repeated-task performance.

Key takeaway:
- Memory is most useful when it captures compact lessons from failures, not giant transcripts.

Implication for agents:
- Store short, actionable reflections.
- Use memory across repeated tasks or sessions, not as an excuse to keep every token forever.

## System Design And Orchestration

### 9. The Shift from Models to Compound AI Systems (BAIR, 2024)
Source:
- https://bair.berkeley.edu/blog/2024/02/18/compound-ai-systems/

Why it matters:
- Strong AI systems increasingly come from multiple interacting components, not just bigger base models.

Key takeaways from the article:
- system design can improve quality faster than scaling alone
- current-data access, control, trust, and cost are often easier to solve at the system level
- optimizing a compound system is a distinct engineering problem

Implication for agents:
- Build around tools, retrievers, graders, and routers when they solve a real product problem.
- Do not mistake "compound system" for "maximally complex system".

### 10. Building Effective AI Agents (Anthropic, 2026)
Source:
- https://resources.anthropic.com/building-effective-ai-agents

Why it matters:
- High-quality practical guidance from a team operating real agent systems at scale.

The most useful framing:
- choose between single-agent, workflow, and multi-agent designs intentionally
- use a small set of reusable patterns:
  - sequential
  - parallel
  - evaluator-optimizer
- match system complexity to business value

Implication for agents:
- Default to simple workflows first.
- Reach for evaluator-optimizer when correctness matters and the task benefits from revision.
- Reach for multi-agent only after single-agent/workflow baselines are exhausted.

### 11. Understanding Agent Scaling via Diversity (2026)
Source:
- arXiv:2602.03794

Why it matters:
- More homogeneous agents do not scale indefinitely; diversity matters more than count.

Key takeaway:
- Two meaningfully different agents can outperform a swarm of same-ish agents.

Implication for agents:
- Diversity should come from role, model, tool access, or evidence channel.
- Do not duplicate the same model/prompt ten times and call it orchestration.

### 12. SOLVE-Med / MATA / Small-Model Orchestration (2025-2026)
Sources:
- SOLVE-Med: arXiv:2511.03542
- MATA: arXiv:2602.09642

Why they matter:
- Small specialized models, when orchestrated well, can outperform or match much larger standalone systems.

Key takeaway:
- Cheap specialists for mechanical subproblems are a real design pattern, not a hack.

Implication for agents:
- Route grep/read/run/simple classification to cheaper lanes.
- Reserve expensive models for hard reasoning or integration steps.

### 13. Agent READMEs: An Empirical Study of Context Files for Agentic Coding (2025)
Source:
- arXiv:2511.12884

Why it matters:
- Agent context files become living operational artifacts, but often drift into unreadable piles.

Key takeaways from the study:
- teams heavily specify build/run, architecture, and implementation context
- security and performance are badly underspecified

Implication for agents:
- Keep context files short, operational, and constraint-rich.
- Add explicit non-functional requirements.
- Treat agent context as maintained configuration, not lore.

### 14. From Biased Chatbots to Biased Agents (2026)
Source:
- arXiv:2602.12285

Why it matters:
- Persona baggage can actively hurt agent behavior.

Key takeaway:
- Capability framing helps; character acting often hurts.

Implication for agents:
- Keep personalities light.
- Put behavior into constraints and tools, not theatrics.

### 15. Emergent Coordination in Multi-Agent Systems (2025)
Source:
- arXiv:2510.05174

Why it matters:
- Coordination is better when agents share objectives and understand complementary roles.

Key takeaway:
- Role awareness is useful; vague social-role prompts are not enough.

Implication for agents:
- When using multiple agents, explicitly describe what each one contributes and how outputs combine.

## Software Engineering Agents

### 16. SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering (Yang et al., 2024)
Source:
- https://arxiv.org/abs/2405.15793
- https://github.com/princeton-nlp/SWE-agent

Why it matters:
- The interface between model and environment is part of the model's performance.

Key takeaway:
- Good repo agents need a disciplined shell/editor/test interface, not just a better prompt.

Implication for agents:
- Design the action surface carefully.
- Short loops over read/search/edit/test beat abstract planning without execution.

### 17. Agentless: Demystifying LLM-based Software Engineering Agents (Xia et al., 2024)
Source:
- https://arxiv.org/abs/2407.01489

Why it matters:
- A simpler pipeline can outperform complex software agents at lower cost.

Key takeaway:
- Simpler decomposition often beats a giant autonomous loop.

Implication for agents:
- Always benchmark against a simpler non-agentic or lightly agentic baseline.
- If a full agent loop is not clearly better, cut it.

## Evaluation And Reliability

### 18. Demystifying evals for AI agents (Anthropic, 2026)
Source:
- https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents

Why it matters:
- This is one of the best practical writeups on agent evals and reliability.

High-signal takeaways:
- start early; 20-50 tasks is enough to begin
- write unambiguous tasks with reference solutions
- evaluate both "should do X" and "should not do X"
- isolate trials from each other
- grade outputs/outcomes, not rigid exact traces
- calibrate model graders against humans
- read transcripts constantly
- treat eval-driven development as normal engineering

Implication for agents:
- Search/tool-use policies should be evaluated on both over-triggering and under-triggering.
- Coding agents should be graded on artifacts and tests, not whether they followed a specific thought path.

## Projects Worth Studying

These are not all research papers, but they are useful design references.

### 1. karpathy/autoresearch
Source:
- https://github.com/karpathy/autoresearch

What to study:
- extremely narrow loop
- fixed optimization target
- small mutable surface
- experiment-first framing instead of "general agent platform"

Copy:
- tight loop, fixed budget, metric-first automation

Avoid:
- generalizing it into a broad orchestration layer unless evals justify it

### 2. davebcn87/pi-autoresearch
Source:
- https://github.com/davebcn87/pi-autoresearch

What to study:
- practical extension of the autoresearch idea
- explicit session files
- checks vs crashes vs metric logs
- dashboard and widget feedback

Copy:
- make experiment state visible
- distinguish correctness failures from benchmark failures
- commit only after the right checks pass

### 3. SWE-agent / mini-SWE-agent
Source:
- https://github.com/princeton-nlp/SWE-agent

What to study:
- repo-focused action surface
- issue -> inspect -> edit -> test loop
- benchmark-first iteration

Copy:
- narrow interface and strong harnessing

### 4. OpenHands
Source:
- https://github.com/All-Hands-AI/OpenHands

What to study:
- broad workspace/runtime architecture
- interactive software agent product design

Copy carefully:
- runtime ergonomics and environment handling

Risk:
- very easy to absorb too much framework complexity

### 5. aider's architect/editor split
Source:
- https://aider.chat/2024/09/26/architect.html

What to study:
- separate high-level reasoning from concrete editing

Copy:
- planner/editor separation can help when one lane should stay terse and execution-oriented

Risk:
- only worth it if the split clearly improves results on your tasks

## Distilled Rules For Kokoclaw/OpenClaw-Like Systems

### Search And Retrieval

- Do not rely on hardcoded keywords to decide whether to search.
- Let the model judge whether fresh evidence is needed, then measure the behavior.
- Keep the first search shallow and literal.
- Allow bounded refinement if the first results are weak or mismatched.
- Ground final factual answers in retrieved evidence.

### Coding Agents

- Keep repo agents on short inspect/edit/test loops.
- Preserve exact names and file-local conventions.
- Use docs lookup when behavior depends on framework or version details.
- Grade the produced diff and test result, not exact intermediate steps.
- Always compare against simpler baselines like "read more, act less".

### Multi-Agent Design

- Use one agent unless there is a measured reason to split.
- Split by capability, not by story or persona.
- Small helpers should do mechanical work.
- Larger models should handle synthesis and edge-case reasoning.
- Coordination prompts should name the shared objective and each role's responsibility.

### Prompt Writing

- Short beats bloated.
- Abstract rules beat example catalogs unless the task genuinely needs demonstrations.
- Put uncertainty policy in the prompt:
  - verify
  - revise
  - say unknown when unsupported
- Do not try to encode every failure mode in one mega-prompt.

### Evaluation

- Run the same harness the product actually uses.
- Keep trials isolated.
- Track pass@1 and consistency, not just "found a good answer once".
- Review transcripts every week if the system matters.
- If the model improved but the score did not, suspect the benchmark or grader too.

## Anti-Patterns

- Too many homogeneous agents
- Persona-rich prompts with weak task constraints
- Keyword-triggered search/tool policies
- Unbounded self-reflection loops
- Auto-commits without validation
- Massive context files with no ownership
- Grading only the exact path instead of the delivered outcome
- Building a platform before validating a narrow workflow

## What To Re-Read Often

- Anthropic, Building Effective AI Agents
- Anthropic, Demystifying evals for AI agents
- BAIR, The Shift from Models to Compound AI Systems
- ReAct
- CRITIC
- Agentless
- SWE-agent
- karpathy/autoresearch
- pi-autoresearch

## Update Policy For This File

When adding a new source, prefer one of:
- primary paper
- official engineering article
- official project README or documentation

For each new source, capture:
- what it claims
- what to copy
- what to avoid
- whether it actually changes system design decisions

If it does not change design decisions, it probably does not belong here.

*Last updated: 2026-03-25*
