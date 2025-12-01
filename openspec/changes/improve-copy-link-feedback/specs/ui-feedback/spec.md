## ADDED Requirements

### Requirement: Enhanced Copy Link Feedback
The UI SHALL provide clear and noticeable visual feedback when a user copies an object access link to clipboard.

#### Scenario: Successful copy with prominent feedback
- **WHEN** a user clicks the "获取链接" (Get Link) button on an object
- **THEN** the system copies the URL to clipboard
- **AND** displays a prominent, easily noticeable success notification
- **AND** the notification persists long enough for the user to acknowledge it

#### Scenario: Copy feedback visibility
- **WHEN** the copy success notification is displayed
- **THEN** it is visually distinct from regular messages
- **AND** it does not blend into the background or get missed by users

#### Scenario: Failed copy with clear error
- **WHEN** clipboard copy fails for any reason
- **THEN** the system displays a clear error message
- **OR** falls back to displaying the URL for manual copying
- **AND** the user understands what went wrong

### Requirement: Copy Feedback Persistence
The copy success feedback SHALL remain visible long enough for users to perceive and acknowledge the action.

#### Scenario: Notification duration
- **WHEN** a copy success notification is shown
- **THEN** it remains visible for at least 2-3 seconds
- **OR** can be manually dismissed by the user

#### Scenario: Multiple copy actions
- **WHEN** a user copies multiple links in quick succession
- **THEN** each copy action provides clear feedback
- **AND** notifications do not overlap or confuse the user

### Requirement: Fallback Copy Mechanism Preservation
The system SHALL maintain the existing fallback mechanisms for clipboard operations.

#### Scenario: Modern browser clipboard API
- **WHEN** the browser supports navigator.clipboard API
- **THEN** use it as the primary copy method
- **AND** provide enhanced feedback on success

#### Scenario: Legacy browser fallback
- **WHEN** clipboard API is unavailable
- **THEN** fall back to document.execCommand('copy')
- **AND** still provide clear feedback

#### Scenario: All methods fail
- **WHEN** both clipboard API and execCommand fail
- **THEN** display the URL in a notification or modal
- **AND** inform the user to copy manually
